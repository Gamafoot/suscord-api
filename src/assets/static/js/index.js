
function discordApp() {
    return {
        servers: [{ id: 1, name: "Suscord" }],
        activeServer: 1,
        chats: [],
        activeChat: null,
        chatMessages: {},
        chatStates: {},
        userCache: {},
        currentUser: null,
        onlineMembers: [],
        newMessage: '',
        ws: null,
        isConnected: false,
        loadingMessages: false,
        MESSAGE_WINDOW_SIZE: 200,
        selectedFiles: [],
        reconnectTimer: null,
        dragCounter: 0,

        // WebRTC
        peerConnection: null,
        localStream: null,
        isCallActive: false,
        isMuted: false,
        callStatus: 'Соединение...',
        incomingCall: { show: false, from: '', chatId: null, offer: null, timeLeft: 10 },
        callMembers: [],
        callAudioElsByUserId: new Map(),
        callRemoteStreamsByUserId: new Map(),
        _audioMountedListenerAdded: false,
        remoteVolume: 100,
        noiseSuppression: true,
        inputSensitivity: 0,
        audioContext: null,
        gainNode: null,
        remoteMuted: false,
        callOfferTimer: null,
        incomingCallTimer: null,

        rtcConfig: {
            iceServers: [
                { urls: 'stun:stun.l.google.com:19302' },
                { urls: 'stun:stun1.l.google.com:19302' },
                {
                    urls: 'turn:openrelay.metered.ca:80',
                    username: 'openrelayproject',
                    credential: 'openrelayproject'
                },
                {
                    urls: 'turn:openrelay.metered.ca:443',
                    username: 'openrelayproject',
                    credential: 'openrelayproject'
                }
            ]
        },

        notification: { show: false, message: '', icon: 'ℹ️' },
        currentImageSrc: '',
        showInviteModal: false,
        nonMembers: [],
        inviteStates: {},
        inviteNotification: { show: false, code: '', timeLeft: 10 },
        searchQuery: '',
        searchTimeout: null,
        activeTab: 'chats',
        searchResults: [],
        allChats: [],
        showEditModal: false,
        editChatName: '',
        selectedAvatar: null,
        avatarPreview: null,

        init() {
            window.servers = this.servers;
            this.loadCurrentUser();
            this.loadChats().then(() => {
                const pathParts = window.location.pathname.split('/');
                if (pathParts[1] === 'chats' && pathParts[2]) {
                    const chatId = parseInt(pathParts[2]);
                    if (chatId) {
                        this.activeChat = null;
                        this.selectChat(chatId);
                    }
                }
            });
            this.connectWebSocket();

            // Очищаем поиск при переключении вкладок
            this.$watch('activeTab', () => {
                this.searchQuery = '';
                this.searchResults = [];
            });

            window.addEventListener('popstate', () => {
                const pathParts = window.location.pathname.split('/');
                if (pathParts[1] === 'chats' && pathParts[2]) {
                    const chatId = parseInt(pathParts[2]);
                    if (chatId) {
                        this.activeChat = null;
                        this.selectChat(chatId);
                    }
                } else {
                    this.activeChat = null;
                }
            });
        },

        async loadChats() {
            try {
                const res = await fetch('/api/v1/chats');
                if (res.ok) {
                    this.allChats = [...(await res.json())];
                    this.chats = [...this.allChats];
                    return true;
                }
            } catch (e) { console.error(e); }
            return false;
        },

        get messages() {
            return this.chatMessages[this.activeChat] || [];
        },

        async selectChat(chatId) {
            if (this.activeChat === chatId) return;

            this.activeChat = chatId;
            window.history.pushState({}, '', `/chats/${chatId}`);

            if (!this.chatMessages[chatId]) {
                this.chatMessages[chatId] = [];
                this.chatStates[chatId] = {
                    oldestMessageId: null,
                    hasMoreMessages: true,
                    isInitialLoad: true
                };
                await this.loadLatestMessages();
            }

            await this.loadChatMembers();
            setTimeout(() => this.scrollToBottom(), 100);
        },

        async loadLatestMessages() {
            if (!this.activeChat) return;
            this.loadingMessages = true;
            try {
                const url = `/api/v1/chats/${this.activeChat}/messages?chat_id=${this.activeChat}&limit=50`;
                const res = await fetch(url);
                if (res.ok) {
                    const msgs = await res.json();
                    if (msgs.length > 0) {
                        const processedMsgs = await this.processMessages(msgs.reverse());
                        this.chatMessages[this.activeChat] = processedMsgs;
                        this.chatStates[this.activeChat].oldestMessageId = processedMsgs[0].id;
                        this.chatStates[this.activeChat].hasMoreMessages = msgs.length === 50;
                        this.chatStates[this.activeChat].isInitialLoad = false;
                    }
                }
            } catch (e) { console.error(e); }
            this.loadingMessages = false;
        },

        async loadOlderMessages() {
            const state = this.chatStates[this.activeChat];
            if (!this.activeChat || !state?.hasMoreMessages || this.loadingMessages) return;

            this.loadingMessages = true;
            const container = document.getElementById('messages-container');
            const oldScrollHeight = container.scrollHeight;

            try {
                let url = `/api/v1/chats/${this.activeChat}/messages?limit=50`;
                if (state.oldestMessageId) url += `&last_message_id=${state.oldestMessageId}`;

                const res = await fetch(url);
                if (res.ok) {
                    const msgs = await res.json();
                    if (msgs.length > 0) {
                        const processedMsgs = await this.processMessages(msgs.reverse());
                        state.oldestMessageId = processedMsgs[0].id;
                        state.hasMoreMessages = msgs.length === 50;

                        this.chatMessages[this.activeChat] = [...processedMsgs, ...this.chatMessages[this.activeChat]];
                        this.trimMessages();

                        setTimeout(() => {
                            const newScrollHeight = container.scrollHeight;
                            container.scrollTop = newScrollHeight - oldScrollHeight;
                        }, 0);
                    } else {
                        state.hasMoreMessages = false;
                    }
                }
            } catch (e) { console.error(e); }
            this.loadingMessages = false;
        },

        trimMessages() {
            const messages = this.chatMessages[this.activeChat];
            if (messages && messages.length > this.MESSAGE_WINDOW_SIZE) {
                this.chatMessages[this.activeChat] = messages.slice(-this.MESSAGE_WINDOW_SIZE);
                this.chatStates[this.activeChat].hasMoreMessages = true;
            }
        },

        onScroll(e) {
            if (e.target.scrollTop < 100) {
                this.loadOlderMessages();
            }
        },

        scrollToBottom() {
            const container = document.getElementById('messages-container');
            if (container) container.scrollTop = container.scrollHeight;
        },

        async sendMessage() {
            if ((!this.newMessage.trim() && this.selectedFiles.length === 0) || !this.activeChat) return;

            const content = this.newMessage;
            const files = [...this.selectedFiles];
            this.newMessage = '';
            this.selectedFiles = [];

            try {
                const formData = new FormData();

                // Определяем тип сообщения и контент
                if (files.length > 0) {
                    formData.append('type', 'files');
                    if (content.trim()) {
                        formData.append('content', content);
                    }
                } else {
                    formData.append('type', 'message');
                    formData.append('content', content);
                }

                for (const file of files) {
                    formData.append('file', file);
                }

                const res = await fetch(`/api/v1/chats/${this.activeChat}/messages`, {
                    method: 'POST',
                    body: formData
                });

                if (!res.ok) {
                    this.showNotification('Ошибка отправки сообщения', '❌');
                    this.newMessage = content;
                    this.selectedFiles = files;
                    return;
                }

                this.moveChatToTop(this.activeChat);
                setTimeout(() => this.scrollToBottom(), 100);
            } catch (e) {
                console.error('Send message error:', e);
                this.showNotification('Ошибка отправки сообщения', '❌');
                this.newMessage = content;
                this.selectedFiles = files;
            }
        },

        handleFileSelect(event) {
            const files = Array.from(event.target.files);
            this.selectedFiles = [...this.selectedFiles, ...files];
            event.target.value = '';
        },

        handleFileDrop(event) {
            event.preventDefault();
            this.dragCounter = 0;
            const files = Array.from(event.dataTransfer.files);
            this.selectedFiles = [...this.selectedFiles, ...files];
        },

        handleDragEnter() {
            this.dragCounter++;
        },

        handleDragLeave() {
            this.dragCounter--;
        },

        removeFile(index) {
            this.selectedFiles.splice(index, 1);
        },

        isImageOrVideo(file) {
            return file.type.startsWith('image/') || file.type.startsWith('video/');
        },

        isGif(mimeType) {
            return mimeType === 'image/gif';
        },

        getFilePreview(file) {
            if (file.type.startsWith('image/')) {
                return URL.createObjectURL(file);
            }
            return '';
        },

        formatFileSize(bytes) {
            if (bytes === 0) return '0 Б';
            const k = 1024;
            const sizes = ['Б', 'КБ', 'МБ', 'ГБ'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
        },

        isMediaFile(mimeType) {
            return mimeType && (mimeType.startsWith('image/') || mimeType.startsWith('video/') || mimeType.startsWith('audio/'));
        },

        downloadFile(filePath) {
            const link = document.createElement('a');
            link.href = filePath;
            link.download = filePath.split('/').pop();
            link.click();
        },

        async deleteAttachment(attachmentId) {
            if (!confirm('Вы уверены, что хотите удалить этот файл?')) return;

            const { ok } = await this.apiCall(`/api/v1/attachments/${attachmentId}`, { method: 'DELETE' });
            if (ok) {
                const messages = this.chatMessages[this.activeChat];
                if (messages) {
                    for (let i = messages.length - 1; i >= 0; i--) {
                        const message = messages[i];
                        if (message.attachments) {
                            message.attachments = message.attachments.filter(att => att.id !== attachmentId);
                            if (message.attachments.length === 0 && !message.content?.trim()) {
                                messages.splice(i, 1);
                            }
                        }
                    }
                }
                this.showNotification('Файл удалён', '✅');
            } else {
                this.showNotification('Ошибка удаления', '❌');
            }
        },

        getActiveChat() {
            return this.allChats.find(c => c.id === this.activeChat) || this.chats.find(c => c.id === this.activeChat);
        },

        getActiveChatName() {
            const chat = this.getActiveChat();
            return chat ? chat.name : '';
        },

        getActiveFriendName() {
            const chat = this.getActiveChat();
            return chat ? chat.name : 'Пользователь';
        },

        async createNewChat() {
            try {
                const name = prompt('Название группы:');
                if (!name) return;

                const res = await fetch('/api/v1/chats/group', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        name: name,
                        avatar_url: ''
                    })
                });

                if (res.ok) {
                    const chat = await res.json();
                    this.chats.unshift(chat);
                    this.showNotification('Чат создан', '✅');
                } else {
                    this.showNotification('Ошибка создания чата', '❌');
                }
            } catch (e) {
                console.error(e);
                this.showNotification('Ошибка создания чата', '❌');
            }
        },

        async inviteUserToChat() {
            const chat = this.getActiveChat();
            if (!chat || chat.type !== 'group') {
                this.showNotification('Можно приглашать только в групповые чаты', '❌');
                return;
            }

            await this.loadNonMembers();
            this.showInviteModal = true;
        },

        async loadNonMembers() {
            try {
                const res = await fetch(`/api/v1/chats/${this.activeChat}/non-members`);
                if (res.ok) {
                    this.nonMembers = await res.json();
                    this.inviteStates = {};
                }
            } catch (e) {
                console.error(e);
            }
        },

        async sendInvite(userId) {
            try {
                const res = await fetch(`/api/v1/chats/${this.activeChat}/invite`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ user_id: userId })
                });

                if (res.ok) {
                    this.inviteStates[userId] = true;
                } else {
                    this.showNotification('Ошибка отправки приглашения', '❌');
                }
            } catch (e) {
                console.error(e);
                this.showNotification('Ошибка отправки приглашения', '❌');
            }
        },

        async loadChatMembers() {
            if (!this.activeChat) return;

            try {
                const res = await fetch(`/api/v1/chats/${this.activeChat}/members`);
                if (res.ok) {
                    this.onlineMembers = await res.json();
                }
            } catch (e) {
                console.error(e);
            }
        },

        showInviteNotification(code) {
            this.inviteNotification = { show: true, code: code, timeLeft: 10 };

            const timer = setInterval(() => {
                this.inviteNotification.timeLeft--;
                if (this.inviteNotification.timeLeft <= 0) {
                    this.inviteNotification.show = false;
                    clearInterval(timer);
                }
            }, 1000);
        },

        async acceptInvite() {
            try {
                const res = await fetch(`/api/v1/chats/invite/accept/${this.inviteNotification.code}`);
                if (res.ok) {
                    this.inviteNotification.show = false;
                    await this.loadChats();
                    this.showNotification('Приглашение принято', '✅');
                } else {
                    this.showNotification('Ошибка принятия приглашения', '❌');
                }
            } catch (e) {
                console.error(e);
                this.showNotification('Ошибка принятия приглашения', '❌');
            }
        },

        declineInvite() {
            this.inviteNotification.show = false;
        },

        async leaveChat() {
            if (!this.activeChat || !confirm('Вы уверены, что хотите покинуть этот чат?')) return;

            try {
                const res = await this.apiCall(`/api/v1/chats/${this.activeChat}/leave`);
                if (res.ok) {
                    this.showNotification('Вы покинули чат', '✓');
                    this.activeChat = null;
                    await this.loadChats();
                }
            } catch (e) {
                console.error(e);
                this.showNotification('Ошибка при выходе из чата', '⚠️');
            }
        },

        async deleteChat() {
            if (!this.activeChat || !confirm('Вы уверены, что хотите удалить этот чат?')) return;

            try {
                const res = await this.apiCall(`/api/v1/chats/${this.activeChat}`, { method: 'DELETE' });
                if (res.ok) {
                    this.showNotification('Чат удалён', '✅');
                    this.activeChat = null;
                    await this.loadChats();
                } else {
                    this.showNotification('Ошибка удаления чата', '❌');
                }
            } catch (e) {
                console.error(e);
                this.showNotification('Ошибка удаления чата', '❌');
            }
        },

        editGroupChat() {
            const chat = this.getActiveChat();
            if (!chat || chat.type !== 'group') return;

            this.editChatName = chat.name;
            this.selectedAvatar = null;
            this.avatarPreview = null;
            this.showEditModal = true;
        },

        cropTool() {
            return {
                cropArea: { x: 50, y: 50, width: 200, height: 200 },
                isDragging: false,
                isResizing: false,
                startPos: { x: 0, y: 0 },
                resizeStartArea: null,
                dragHandler: null,
                stopHandler: null,
                resizeHandler: null,
                showCrop: true,

                initCrop() {
                    const img = this.$refs.cropImage;
                    if (!img) return;

                    // Проверяем оригинальный размер изображения
                    if (img.naturalWidth < 300 || img.naturalHeight < 300) {
                        this.showCrop = false;
                        return;
                    }

                    this.showCrop = true;

                    // Ждем пока изображение отрендерится
                    const initCropArea = () => {
                        const rect = img.getBoundingClientRect();
                        if (rect.width === 0 || rect.height === 0) {
                            requestAnimationFrame(initCropArea);
                            return;
                        }
                        const size = Math.min(rect.width, rect.height, 600);
                        this.cropArea = {
                            x: (rect.width - size) / 2,
                            y: (rect.height - size) / 2,
                            width: size,
                            height: size
                        };
                    };
                    requestAnimationFrame(initCropArea);
                },

                startDrag(e) {
                    this.isDragging = true;
                    const img = this.$refs.cropImage;
                    const rect = img.getBoundingClientRect();
                    this.startPos = {
                        x: e.clientX - rect.left - this.cropArea.x,
                        y: e.clientY - rect.top - this.cropArea.y
                    };
                    this.dragHandler = this.drag.bind(this);
                    this.stopHandler = this.stopDrag.bind(this);
                    document.addEventListener('mousemove', this.dragHandler);
                    document.addEventListener('mouseup', this.stopHandler);
                },

                drag(e) {
                    if (!this.isDragging) return;
                    const img = this.$refs.cropImage;
                    const rect = img.getBoundingClientRect();
                    const newX = Math.max(0, Math.min(e.clientX - rect.left - this.startPos.x, img.offsetWidth - this.cropArea.width));
                    const newY = Math.max(0, Math.min(e.clientY - rect.top - this.startPos.y, img.offsetHeight - this.cropArea.height));
                    this.cropArea.x = newX;
                    this.cropArea.y = newY;
                },

                stopDrag() {
                    this.isDragging = false;
                    if (this.dragHandler) document.removeEventListener('mousemove', this.dragHandler);
                    if (this.stopHandler) document.removeEventListener('mouseup', this.stopHandler);
                },

                startResize(e, corner) {
                    e.preventDefault();
                    this.isResizing = true;
                    const img = this.$refs.cropImage;
                    const rect = img.getBoundingClientRect();
                    this.startPos = { x: e.clientX, y: e.clientY };
                    this.resizeStartArea = { ...this.cropArea };
                    this.resizeCorner = corner;
                    this.resizeHandler = this.resize.bind(this);
                    this.stopHandler = this.stopResize.bind(this);
                    document.addEventListener('mousemove', this.resizeHandler);
                    document.addEventListener('mouseup', this.stopHandler);
                },

                resize(e) {
                    if (!this.isResizing) return;
                    const img = this.$refs.cropImage;
                    const dx = e.clientX - this.startPos.x;
                    const dy = e.clientY - this.startPos.y;
                    const delta = Math.max(dx, dy);

                    let newSize = this.resizeStartArea.width + delta;
                    let newX = this.resizeStartArea.x;
                    let newY = this.resizeStartArea.y;

                    if (this.resizeCorner === 'nw') {
                        newSize = this.resizeStartArea.width - delta;
                        newX = this.resizeStartArea.x + delta;
                        newY = this.resizeStartArea.y + delta;
                    } else if (this.resizeCorner === 'ne') {
                        newSize = this.resizeStartArea.width + delta;
                        newY = this.resizeStartArea.y - delta;
                    } else if (this.resizeCorner === 'sw') {
                        newSize = this.resizeStartArea.width - delta;
                        newX = this.resizeStartArea.x + delta;
                    }

                    newSize = Math.max(50, Math.min(newSize, img.offsetWidth, img.offsetHeight, 600));
                    newX = Math.max(0, Math.min(newX, img.offsetWidth - newSize));
                    newY = Math.max(0, Math.min(newY, img.offsetHeight - newSize));

                    this.cropArea = { x: newX, y: newY, width: newSize, height: newSize };
                },

                stopResize() {
                    this.isResizing = false;
                    if (this.resizeHandler) document.removeEventListener('mousemove', this.resizeHandler);
                    if (this.stopHandler) document.removeEventListener('mouseup', this.stopHandler);
                }
            };
        },

        handleAvatarSelect(event) {
            const file = event.target.files[0];
            if (file) {
                const reader = new FileReader();
                reader.onload = (e) => {
                    const img = new Image();
                    img.onload = () => {
                        if (img.width < 300 || img.height < 300) {
                            this.showNotification('Размер изображения должен быть минимум 300x300 пикселей', '❌');
                            event.target.value = '';
                            return;
                        }
                        this.selectedAvatar = file;
                        this.avatarPreview = e.target.result;
                    };
                    img.src = e.target.result;
                };
                reader.readAsDataURL(file);
            }
        },

        async getCroppedImage() {
            if (!this.avatarPreview) return this.selectedAvatar;

            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');
            const img = new Image();

            return new Promise((resolve) => {
                img.onload = () => {
                    const cropOverlay = document.querySelector('.crop-overlay');
                    const cropImg = document.querySelector('.crop-image');

                    // Если изображение маленькое или нет overlay, возвращаем оригинал
                    if (!cropOverlay || !cropImg || img.naturalWidth < 300 || img.naturalHeight < 300) {
                        resolve(this.selectedAvatar);
                        return;
                    }

                    const rect = cropOverlay.getBoundingClientRect();
                    const imgRect = cropImg.getBoundingClientRect();

                    const scaleX = img.naturalWidth / imgRect.width;
                    const scaleY = img.naturalHeight / imgRect.height;

                    const cropX = (rect.left - imgRect.left) * scaleX;
                    const cropY = (rect.top - imgRect.top) * scaleY;
                    const cropWidth = rect.width * scaleX;
                    const cropHeight = rect.height * scaleY;

                    // Ограничиваем размер итогового изображения до 600x600
                    const maxSize = 600;
                    const finalWidth = Math.min(cropWidth, maxSize);
                    const finalHeight = Math.min(cropHeight, maxSize);

                    canvas.width = finalWidth;
                    canvas.height = finalHeight;

                    ctx.drawImage(img, cropX, cropY, cropWidth, cropHeight, 0, 0, finalWidth, finalHeight);

                    canvas.toBlob((blob) => {
                        const croppedFile = new File([blob], this.selectedAvatar.name, { type: this.selectedAvatar.type });
                        resolve(croppedFile);
                    }, this.selectedAvatar.type);
                };
                img.src = this.avatarPreview;
            });
        },

        async saveGroupChat() {
            if (!this.activeChat) return;

            try {
                const formData = new FormData();

                if (this.editChatName.trim()) {
                    formData.append('name', this.editChatName.trim());
                }

                if (this.selectedAvatar) {
                    const croppedFile = await this.getCroppedImage();
                    formData.append('file', croppedFile);
                }

                const res = await fetch(`/api/v1/chats/${this.activeChat}`, {
                    method: 'PATCH',
                    body: formData
                });

                if (res.ok) {
                    const updatedChat = await res.json();

                    // Обновляем чат в списке
                    const chatIndex = this.chats.findIndex(c => c.id === this.activeChat);
                    if (chatIndex !== -1) {
                        this.chats[chatIndex] = { ...this.chats[chatIndex], ...updatedChat };
                    }

                    const allChatIndex = this.allChats.findIndex(c => c.id === this.activeChat);
                    if (allChatIndex !== -1) {
                        this.allChats[allChatIndex] = { ...this.allChats[allChatIndex], ...updatedChat };
                    }

                    // Сбрасываем переменные
                    this.avatarPreview = null;
                    this.selectedAvatar = null;
                    this.showEditModal = false;
                    this.showNotification('Чат обновлён', '✅');
                } else {
                    this.showNotification('Ошибка обновления чата', '❌');
                }
            } catch (e) {
                console.error(e);
                this.showNotification('Ошибка обновления чата', '❌');
            }
        },

        getLastMessage(chatId) {
            // Заглушка, можно реализовать кэширование последних сообщений
            return '';
        },

        formatTime(ts) {
            if (!ts) return '';
            const date = new Date(ts);
            return date.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' }) + ' ' + date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
        },

        // WebSocket
        connectWebSocket() {
            const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
            const urlWebsocket = `${protocol}//${location.host}/ws?session=${this.getCookie('session')}`;
            this.ws = new WebSocket(urlWebsocket);

            this.ws.onopen = () => {
                this.isConnected = true;
                this.showNotification('Подключено', '✅');

                if (this.reconnectTimer) {
                    clearTimeout(this.reconnectTimer);
                    this.reconnectTimer = null;
                }

                this.hideLoadingScreen();
            };

            this.ws.onclose = () => {
                this.isConnected = false;
                this.showNotification('Отключено от сервера', '🔴');

                if (this.reconnectTimer) {
                    clearTimeout(this.reconnectTimer);
                }

                this.reconnectTimer = setTimeout(() => {
                    this.showLoadingScreen();
                }, 5000);

                setTimeout(() => this.connectWebSocket(), 3000);
            };

            this.ws.onmessage = (event) => {
                const msg = JSON.parse(event.data);
                this.handleWsMessage(msg);
            };
        },

        async handleWsMessage(msg) {
            if (msg.type.includes("call-")) {
                this.handleWebRTCSignaling(msg);
                return;
            }

            switch (msg.type) {
                case 'message':
                    let chatId = msg.data.chat_id;
                    if (!this.chatMessages[chatId]) {
                        this.chatMessages[chatId] = [];
                        this.chatStates[chatId] = {
                            oldestMessageId: null,
                            hasMoreMessages: true,
                            isInitialLoad: true
                        };
                    }

                    // Обрабатываем новую структуру сообщения
                    const messageData = {
                        id: msg.data.id,
                        user_id: msg.data.user_id,
                        content: msg.data.content,
                        type: msg.data.type,
                        timestamp: msg.data.created_at,
                        attachments: msg.data.attachments || []
                    };

                    const processedMsg = await this.processMessage(messageData, true);
                    this.chatMessages[chatId].push(processedMsg);
                    this.trimMessages();

                    this.moveChatToTop(chatId);
                    if (chatId === this.activeChat) {
                        setTimeout(() => this.scrollToBottom(), 100);
                    }
                    break;

                case 'message_update':
                    this.handleMessageUpdate(msg.data);
                    break;

                case 'message_delete':
                    this.handleMessageDelete(msg.data);
                    break;

                case 'joined_chat':
                    {
                        const chat = this.chats.find(c => c.id === msg.data.chat.id);
                        if (!chat) {
                            this.chats.unshift({
                                id: msg.data.chat.id,
                                name: msg.data.chat.name,
                                avatar_url: msg.data.chat.avatar_url,
                                type: 'private'
                            });
                        }
                    }

                    break;

                case 'new_user_in_chat':
                    const member = this.onlineMembers.find(m => m.id === msg.data.user.id);

                    if (msg.data.chat_id === this.activeChat && !member) {
                        this.onlineMembers.push(msg.data.user);
                        this.showNotification(`${msg.data.user.username} зашел в чат`, '👋');
                    }

                    break;

                case 'user_left':
                    if (msg.data.chat_id === this.activeChat) {
                        const member = this.onlineMembers.find(c => c.id === msg.data.user_id);
                        this.showNotification(`${member.username} покинул чат`, '🚪');
                        this.onlineMembers = this.onlineMembers.filter(m => m.id !== msg.data.user_id);
                    }
                    break;

                case 'delete_chat':
                    {
                        if (msg.data.chat_id === this.activeChat) {
                            this.onlineMembers = [];
                            this.activeChat = null;
                        }

                        const chat = this.chats.find(c => c.id === msg.data.chat_id);
                        this.chats = this.chats.filter(c => c.id !== msg.data.chat_id);
                        this.showNotification(`Чат ${chat?.name || '"нет имени"'} удален`, '🗑️');
                    }
                    break;

                case 'invite_to_chat':
                    this.showInviteNotification(msg.data.code);
                    break;

                case 'update_group_chat':
                    {
                        const updatedChat = msg.data.chat;
                        const chatIndex = this.chats.findIndex(c => c.id === updatedChat.id);
                        if (chatIndex !== -1) {
                            this.chats[chatIndex] = { ...this.chats[chatIndex], ...updatedChat };
                        }
                        const allChatIndex = this.allChats.findIndex(c => c.id === updatedChat.id);
                        if (allChatIndex !== -1) {
                            this.allChats[allChatIndex] = { ...this.allChats[allChatIndex], ...updatedChat };
                        }
                        this.showNotification(`Чат ${chat?.name} был обновлен`, '🔄');
                    }
                    break;

                default:
                    console.log('WS:', msg);
            }
        },

        handleMessageUpdate(data) {
            const messages = this.chatMessages[data.chat_id];
            if (messages) {
                const messageIndex = messages.findIndex(m => m.id === data.id);
                if (messageIndex !== -1) {
                    messages[messageIndex] = {
                        ...messages[messageIndex],
                        content: data.content,
                        attachments: data.attachments || []
                    };
                }
            }
        },

        handleMessageDelete(data) {
            this.removeMessage(data.chat_id, data.ID);
        },

        // Утилиты
        async apiCall(url, options = {}) {
            try {
                const res = await fetch(url, options);
                let data = null;
                if (res.ok) {
                    try {
                        data = await res.json();
                    } catch {
                        data = null;
                    }
                }
                return { ok: res.ok, data };
            } catch (e) {
                console.error(e);
                return { ok: false, data: null };
            }
        },

        moveChatToTop(chatId) {
            const index = this.chats.findIndex(c => c.id === chatId);
            if (index > 0) {
                const chat = this.chats.splice(index, 1)[0];
                this.chats.unshift(chat);
            }
        },

        debouncedSearch() {
            clearTimeout(this.searchTimeout);
            this.searchTimeout = setTimeout(() => {
                if (this.activeTab === 'chats') {
                    this.searchChats();
                } else {
                    this.searchUsers();
                }
            }, 300);
        },

        async searchChats() {
            const query = this.searchQuery.trim();
            if (!query) {
                this.chats = [...this.allChats];
                return;
            }
            const { ok, data } = await this.apiCall(`/api/v1/chats?search=${encodeURIComponent(query)}`);
            if (ok && data) this.chats = data;
        },

        async searchUsers() {
            const query = this.searchQuery.trim();
            if (!query) {
                this.searchResults = [];
                return;
            }
            const { ok, data } = await this.apiCall(`/api/v1/users?search=${encodeURIComponent(query)}`);
            if (ok && data) this.searchResults = data;
        },

        async createPrivateChat(userId) {
            try {
                const res = await fetch('/api/v1/chats/private', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ user_id: userId })
                });

                if (res.ok) {
                    const chat = await res.json();
                    const existingChat = this.chats.find(c => c.id === chat.id);
                    if (!existingChat) {
                        this.chats.unshift(chat);
                    }
                    this.activeTab = 'chats';
                    this.selectChat(chat.id);
                    this.searchQuery = '';
                    this.searchResults = [];
                    this.showNotification('Приватный чат создан', '✅');
                } else {
                    this.showNotification('Ошибка создания чата', '❌');
                }
            } catch (e) {
                console.error(e);
                this.showNotification('Ошибка создания чата', '❌');
            }
        },

        findMessage(chatId, messageId) {
            const messages = this.chatMessages[chatId];
            return messages ? messages.findIndex(m => m.id === messageId) : -1;
        },

        removeMessage(chatId, messageId) {
            const index = this.findMessage(chatId, messageId);
            if (index !== -1) this.chatMessages[chatId].splice(index, 1);
        },

        getCookie(name) {
            const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
            return match ? match[2] : '';
        },

        showNotification(message, icon = 'ℹ️') {
            this.notification = { show: true, message, icon };
            setTimeout(() => this.notification.show = false, 4000);
        },

        hideLoadingScreen() {
            const loadingScreen = document.getElementById('loading-screen');
            if (loadingScreen) {
                loadingScreen.classList.add('fade-out');
                setTimeout(() => loadingScreen.remove(), 500);
            }
        },

        showLoadingScreen() {
            if (document.getElementById('loading-screen')) return;
            const loadingScreen = document.createElement('div');
            loadingScreen.id = 'loading-screen';
            loadingScreen.className = 'loading-screen';
            loadingScreen.innerHTML = `
                <div class="loading-spinner"></div>
                <div style="font-size: 1.2em; color: #b9bbbe;">Подключение к серверу...</div>
            `;
            document.body.appendChild(loadingScreen);
        },

        async processMessages(messages) {
            const processed = [];
            for (const msg of messages) {
                processed.push(await this.processMessage(msg));
            }
            return processed;
        },

        async processMessage(msg, isFromWs = false) {
            const user = await this.getUser(msg.user_id);

            return {
                id: msg.id,
                user_id: msg.user_id,
                content: msg.content,
                type: msg.type || 'message',
                timestamp: msg.timestamp || msg.created_at,
                username: user?.username || 'Неизвестный пользователь',
                avatar_url: user?.avatar_url,
                attachments: (msg.attachments || []).map(att => ({
                    id: att.id,
                    file_url: att.file_url,
                    file_size: att.file_size,
                    mime_type: att.mime_type
                }))
            };
        },

        async getUser(userId) {
            if (this.userCache[userId]) {
                return this.userCache[userId];
            }

            try {
                const res = await fetch(`/api/v1/users/${userId}`);
                if (res.ok) {
                    const user = await res.json();
                    this.userCache[userId] = user;
                    return user;
                }
            } catch (e) {
                console.error('Error loading user:', e);
            }
            return null;
        },

        async loadCurrentUser() {
            try {
                const res = await fetch('/api/v1/users/me');
                if (res.ok) {
                    this.currentUser = await res.json();
                }
            } catch (e) {
                console.error('Error loading current user:', e);
            }
        },

        async editMessage(message) {
            const newContent = prompt('Редактировать сообщение:', message.content);
            if (newContent && newContent !== message.content) {
                try {
                    const res = await fetch(`/api/v1/messages/${message.id}`, {
                        method: 'PATCH',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ content: newContent })
                    });
                    if (res.ok) {
                        message.content = newContent;
                        this.showNotification('Сообщение отредактировано', '✅');
                    }
                } catch (e) {
                    console.error('Edit message error:', e);
                    this.showNotification('Ошибка редактирования', '❌');
                }
            }
        },

        async deleteMessage(messageId) {
            if (!confirm('Удалить сообщение?')) return;

            const { ok } = await this.apiCall(`/api/v1/messages/${messageId}`, { method: 'DELETE' });
            if (ok) {
                this.removeMessage(this.activeChat, messageId);
                this.showNotification('Сообщение удалено', '✅');
            } else {
                this.showNotification('Ошибка удаления', '❌');
            }
        },

        showImagePopup(imageSrc) {
            this.currentImageSrc = imageSrc;
            document.getElementById('modalImage').src = imageSrc;
            new bootstrap.Modal(document.getElementById('imageModal')).show();
        },

        tiles: new Map(),
        signal: null,
        client: null,
        localStream: null,

        // WebRTC Звонки (SFU)
        async startCall(chatId) {
            let callingChatId = null;

            if (!this.activeChat && !chatId) return;

            if (this.activeChat) {
                callingChatId = this.activeChat;
            } else if (chatId) {
                callingChatId = chatId;
            }

            this.incomingCall = {
                chatId: callingChatId
            };

            this.isCallActive = true;
            this.callStatus = 'Подлючение...';

            try {
                if (this.currentUser?.id) {
                    const me = {
                        id: this.currentUser.id,
                        username: this.currentUser.username,
                        avatar_url: this.currentUser.avatar_url
                    };
                    this.callMembers = [me];
                }

                if (this.ws && this.ws.readyState === WebSocket.OPEN) {
                    this.ws.send(JSON.stringify({
                        type: 'call-invite',
                        chat_id: callingChatId,
                    }));
                }

            } catch (err) {
                console.error('Ошибка при запуске звонка:', err);
                this.showNotification('Не удалось получить доступ к микрофону', '❌');
                this.endCall();
            }
        },

        upsertCallMember(client) {
            if (!client?.id) return;

            const idx = (this.callMembers || []).findIndex(c => c.id === client.id);
            if (idx === -1) {
                this.callMembers = [...(this.callMembers || []), client];
            } else {
                const next = [...this.callMembers];
                next[idx] = { ...next[idx], ...client };
                this.callMembers = next;
            }
        },

        handleWebRTCSignaling(msg) {
            switch (msg.type) {
                case 'call-invite': {
                    if (this.isCallActive) return;
                    const chatId = msg.chat_id;
                    const chat = this.chats?.find(c => c.id === chatId);

                    this.incomingCall = {
                        show: true,
                        from: chat?.name || 'чат',
                        chatId,
                        offer: null,
                        timeLeft: 10
                    };

                    if (this.incomingCallTimer) clearInterval(this.incomingCallTimer);
                    this.incomingCallTimer = setInterval(() => {
                        if (!this.incomingCall?.show) {
                            clearInterval(this.incomingCallTimer);
                            this.incomingCallTimer = null;
                            return;
                        }

                        this.incomingCall.timeLeft -= 1;
                        if (this.incomingCall.timeLeft <= 0) {
                            clearInterval(this.incomingCallTimer);
                            this.incomingCallTimer = null;
                            this.rejectCall();
                        }
                    }, 1000);
                    break;
                }

                case 'call-accept': {
                    // Кто-то присоединился к звонку
                    this.upsertCallMember(msg.data);

                    this.connectToCall(this.activeChat);

                    // Для инициатора: переключаем UI в "В звонке" после первого присоединения
                    if (this.isCallActive) {
                        this.callStatus = 'В звонке';
                    }

                    console.log("user was accept invite: send stream");
                    this.wsSendStream();
                    break;
                }

                case 'call-reject': {
                    this.endCall({ notifyServer: false });
                    break;
                }

                case 'call-clients': {
                    this.callMembers = msg.data?.clients;
                    break;
                }

                case 'call-leave': {
                    this.callMembers = msg.data?.clients;

                    if (this.callMembers.length === 1) {
                        this.showNotification('Звонок завершён', '📞');
                        this.endCall({ notifyServer: false });
                    }

                    break;
                }

                case 'call-ended': {
                    this.showNotification('Звонок завершён', '📞');
                    this.endCall({ notifyServer: false });
                    break;
                }

                case 'call-stream': {
                    console.log("set stream", msg);
                    this.callUserStreams.set(msg.data.stream_id, msg.data.user_id);
                    break;
                }
            }
        },

        callUserStreams: new Map(),

        connectToCall(chatId) {
            if (this.client) {
                return;
            }

            const protocol = location.protocol.includes("https") ? "wss" : "ws";
            const wsURL = `${protocol}://${location.hostname}:7002/ws`;

            const rtcConfig = {
                iceTransportPolicy: 'all',
                bundlePolicy: 'max-bundle',
            };

            this.signal = new Signal.IonSFUJSONRPCSignal(wsURL);
            this.client = new IonSDK.Client(this.signal, rtcConfig);

            this.client.ontrack = (track, stream) => {
                console.log('track', track);
                console.log('stream', stream);

                const userId = this.callUserStreams.get(stream.id);
                const member = this.callMembers.find(m => m.id === userId);

                if (member) {
                    member.stream = stream;
                }

                track.onended = () => {
                    if (userId != null) {
                        this.removeRemoteAudioByUserId(userId);
                    }
                };
            };

            this.signal.onopen = async () => {
                console.log("on open", '' + chatId, '' + this.currentUser.id);
                this.client.join('' + chatId, '' + this.currentUser.id);

                this.localStream = await IonSDK.LocalStream.getUserMedia({
                    audio: true,
                    video: false,
                    simulcast: false,
                });

                this.wsSendStream();

                this.client.publish(this.localStream);
            };

            this.signal.onerror = (e) => {
                console.error('signal error', e);
                this.showNotification('Ошибка подключения', '⚠️');
                this.incomingCall = {};
                this.callStatus = null;
            };
        },

        wsSendStream() {
            if (this.localStream && this.ws) {
                this.ws.send(JSON.stringify({
                    type: "call-stream",
                    chat_id: this.incomingCall.chatId,
                    data: {
                        user_id: this.currentUser.id,
                        stream_id: this.localStream.id
                    }
                }));
            }
        },

        async acceptCall() {
            if (!this.incomingCall.chatId) return;

            if (this.incomingCallTimer) {
                clearInterval(this.incomingCallTimer);
                this.incomingCallTimer = null;
            }

            const callChatId = this.incomingCall.chatId;
            this.incomingCall.show = false;
            this.isCallActive = true;
            this.callStatus = 'В звонке';

            this.ws.send(JSON.stringify({
                type: 'call-accept',
                chat_id: callChatId,
            }));

            // Просто начинаем звонок (подключаемся к SFU)
            this.connectToCall(callChatId);

            console.log("accept call: send stream");
            this.wsSendStream();
        },

        rejectCall() {
            this.callStatus = null;

            if (this.incomingCallTimer) {
                clearInterval(this.incomingCallTimer);
                this.incomingCallTimer = null;
            }
            if (this.ws && this.ws.readyState === WebSocket.OPEN) {
                this.ws.send(JSON.stringify({
                    type: 'call-reject',
                    chat_id: this.incomingCall.chatId,
                }));
            }
            this.incomingCall = { show: false, from: '', chatId: null, offer: null, timeLeft: 10 };
        },

        endCall({ notifyServer = true } = {}) {
            if (this.callOfferTimer) {
                clearInterval(this.callOfferTimer);
                this.callOfferTimer = null;
            }

            if (this.incomingCallTimer) {
                clearInterval(this.incomingCallTimer);
                this.incomingCallTimer = null;
            }

            if (this.localStream) {
                this.localStream.getTracks().forEach(track => { track.stop(); });
                this.localStream = null;
            }

            if (this.client) {
                try {
                    this.client.leave();
                } catch (err) {
                    console.warn("client.leave failed:", err);
                }
            }

            if (this.client?.pc) {
                this.client.pc.close();
            }

            this.client = null;

            if (this.audioContext) {
                this.audioContext.close();
                this.audioContext = null;
                this.gainNode = null;
            }

            if (notifyServer && this.isCallActive && this.ws && this.ws.readyState === WebSocket.OPEN) {
                this.ws.send(JSON.stringify({
                    type: 'call-leave',
                    chat_id: this.incomingCall.chatId,
                    data: {
                        user_id: this.currentUser.id
                    }
                }));
            }

            // Чистим удалённые аудио/участников
            for (const userId of this.callAudioElsByUserId.keys()) {
                if (userId !== this.currentUser?.id) {
                    this.removeRemoteAudioByUserId(userId);
                }
            }
            this.callMembers = [];
            this.callAudioElsByUserId.clear();
            this.callRemoteStreamsByUserId.clear();

            if (this.incomingCall) {
                this.incomingCall = { show: false, from: '', chatId: null, offer: null, timeLeft: 10, chat_id: null };
            }

            this.isCallActive = false;
            this.isMuted = false;
            this.remoteMuted = false;
            this.callStatus = 'Соединение...';
        },

        toggleMute() {
            if (this.localStream) {
                this.localStream.getAudioTracks().forEach(track => {
                    track.enabled = !track.enabled;
                });
                this.isMuted = !this.isMuted;
            }
        },

        setRemoteVolume(volume) {
            this.remoteVolume = volume;
            const remoteAudio = document.querySelector('audio[x-ref="remoteAudio"]');
            if (remoteAudio && !this.remoteMuted) {
                remoteAudio.volume = Math.min(volume / 100, 2.0);
            }
        },

        toggleRemoteMute() {
            this.remoteMuted = !this.remoteMuted;
            const remoteAudio = document.querySelector('audio[x-ref="remoteAudio"]');
            if (remoteAudio) {
                remoteAudio.volume = this.remoteMuted ? 0 : Math.min(this.remoteVolume / 100, 2.0);
            }
        },

        setInputSensitivity(value) {
            this.inputSensitivity = value;
            if (this.gainNode) {
                const gain = Math.pow(10, value / 20);
                this.gainNode.gain.value = gain;
            }
        },

        async toggleNoiseSuppression() {
            this.noiseSuppression = !this.noiseSuppression;
            if (this.localStream && this.isCallActive && this.client) {
                const oldStream = this.localStream;
                const wasMuted = this.isMuted;
                try {
                    this.localStream = await navigator.mediaDevices.getUserMedia({
                        audio: {
                            echoCancellation: true,
                            noiseSuppression: this.noiseSuppression,
                            autoGainControl: true
                        }
                    });

                    if (wasMuted) {
                        this.localStream.getAudioTracks().forEach(track => track.enabled = false);
                    }

                    oldStream.getTracks().forEach(track => track.stop());
                } catch (error) {
                    console.error('Ошибка переключения шумоподавления:', error);
                    this.localStream = oldStream;
                    this.noiseSuppression = !this.noiseSuppression;
                }
            }
        },
    };
}
