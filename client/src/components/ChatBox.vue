<template>
    <div>
        <el-row v-for="(item,index) in msglist" :key="index">
            <el-col>{{item.username}}:{{item.msg}}</el-col>
        </el-row>
        <el-row>
            <el-col>
                <el-input
                    type="textarea"
                    :autosize="{ minRows: 1, maxRows: 4}"
                    placeholder="请输入内容"
                    v-model="msg">
                </el-input>
            </el-col>
        </el-row>
        <el-row>
            <el-col>
                <el-button @click="onSend">发送</el-button>
            </el-col>
        </el-row>
    </div>
</template>
<script>
import { MessageBox } from "element-ui";
    var ws;
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
                this.msglist.push({'userid': '', 'username': '我', 'msg': this.msg})
                ws.send(JSON.stringify({'type':'msg', 'msg': this.msg}))
                this.msg = '';
            },
            onMessage: function(evt) {
                console.log( "Received Message: " + evt.data);
                var data = JSON.parse(evt.data);
                if (data.type == 'msg') {
                    this.msglist.push({'userid': data.senderid, 'username': data.sendername, 'msg': data.msg})
                }
            }
        },
    }
</script>