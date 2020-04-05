<style scoped>
    .use-flex {
        display: inline-flex;
        align-items: flex-end;
        width: 100%;
    }
    .left {
        flex: auto;
    }
    .right {
        flex: none;
        width: 80px;
    }
</style>
<template>
    <el-container>
        <el-aside width="200px">Aside</el-aside>
        <el-container>
            <el-main ref='main'>
                <el-table :data="msglist" stripe style="width: 100%">
                    <el-table-column prop="datetime" label="消息时间" width="120"/>
                    <el-table-column prop="username" label="发送人" width="120"/>
                    <el-table-column prop="msg" label="内容"/>
                </el-table>
            </el-main>
            <el-footer height="40px">
                <div class='use-flex'>
                    <div class='left'>
                        <el-input ref="input" @keydown.enter.native="onEnter" type="textarea" :autosize="{ minRows: 1, maxRows: 4}" placeholder="请输入内容" v-model="msg" />
                    </div>
                    <div class='right'>
                        <el-button size="small" @click="onSend">发送</el-button>
                    </div>
                </div>
            </el-footer>
        </el-container>
    </el-container>
</template>
<script>
    import { MessageBox } from "element-ui";
    import moment from "moment";
    var ws;
    function nowStr() {
        return moment().format('MM-DD hh:mm:ss');
    }
    export default {
        data() {
            return {
                msg: '',
                username: '',
                msglist: [],
            }
        },
        mounted: function() {
            MessageBox.prompt('请输入用户名', '提示', {}).then(({ value }) => {
                if (value == null || value == "") {
                    value = "noname"
                } else {
                    value = value.trim()
                }
                this.createWs(value)
            }).catch(() => {
                this.createWs('noname')
            });
        },
        methods: {
            createWs(userName) {
                var roomId = this.getUrlVars('roomid');
                if (roomId === undefined) {
                    roomId = new Date().getTime().toString() + "." + this.getRandomInt(100000).toString() + this.getRandomInt(100000).toString() + this.getRandomInt(100000).toString()
                    window.history.pushState(null, null, window.location.href + '?roomid=' + roomId)
                }
                console.log('roomId:', roomId)
                ws = new WebSocket(`ws://127.0.0.1/chat?username=${userName}&roomid=${roomId}`);
                ws.onmessage = this.onMessage;
            },
            getRandomInt :function(max) {
                return Math.floor(Math.random() * Math.floor(max));
            },
            getUrlVars :function(parameter) {
                var vars = {};
                window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m,key,value) {
                    vars[key] = value;
                });
                return vars[parameter];
            },
            addMsg: function(senderid, sendername, msg) {
                this.msglist.push({'userid': senderid, 'username': sendername, 'msg': msg, 'datetime': nowStr()})
            },
            onSend: function() {
                var value = this.msg.trim();
                if (value === "") {
                    return;
                }
                this.addMsg('', '我', value)
                ws.send(JSON.stringify({'type':'msg', 'msg': value}))
                this.msg = '';
                this.$refs.input.focus()
            },
            onEnter: function(e) {
                if (e.ctrlKey) {
                    this.onSend()
                    e.preventDefault()
                }
            },
            onMessage: function(evt) {
                console.log( "Received Message: " + evt.data);
                var data = JSON.parse(evt.data);
                switch (data.type) {
                    case 'msg':
                        this.addMsg(data.senderid, data.sendername, data.msg);
                        break
                    case 'join':
                        this.addMsg(data.senderid, '', `${data.sendername}(${data.senderid})加入聊天室`)
                        break
                    case 'leave':
                        this.addMsg(data.senderid, '', `${data.sendername}(${data.senderid})离开聊天室`)
                        break
                }
            }
        },
    }
</script>