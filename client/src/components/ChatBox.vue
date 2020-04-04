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
            MessageBox.prompt('请输入用户名', '提示', {
            }).then(({ value }) => {
                if (value == null || value == "") {
                    value = "noname"
                } else {
                    value = value.trim()
                }
                ws = new WebSocket("ws://127.0.0.1/chat?username="+value+"&roomid=def");
                ws.onmessage = this.onMessage;
            }).catch(() => {
                var value = "noname"
                ws = new WebSocket("ws://127.0.0.1/chat?username="+value+"&roomid=def");
                ws.onmessage = this.onMessage;
            });
        },
        methods: {
            onSend: function() {
                this.msglist.unshift({'userid': '', 'username': '我', 'msg': this.msg, 'datetime': nowStr()})
                ws.send(JSON.stringify({'type':'msg', 'msg': this.msg}))
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
                        this.msglist.unshift({'userid': data.senderid, 'username': data.sendername, 'msg': data.msg, 'datetime': nowStr()})
                        break
                    case 'join':
                        this.msglist.unshift({'userid': data.senderid, 'username': '', 'msg': `${data.sendername}(${data.senderid})加入聊天室`, 'datetime': nowStr()})
                        break
                    case 'leave':
                        this.msglist.unshift({'userid': data.senderid, 'username': '', 'msg': `${data.sendername}(${data.senderid})离开聊天室`, 'datetime': nowStr()})
                        break
                }
            }
        },
    }
</script>