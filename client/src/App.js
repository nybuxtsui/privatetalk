import React from 'react';
import {
  Button,
  Layout,
  Input,
  Table,
  Modal,
  message,
} from 'antd';
import './App.css';
import dayjs from 'dayjs'

const { Header, Footer, Sider, Content } = Layout;
const { TextArea } = Input;


const columns = [
  {
    title: '消息时间',
    dataIndex: 'datetime',
    key: 'datetime',
  },
  {
    title: '发送人',
    dataIndex: 'username',
    key: 'username',
  },
  {
    title: '内容',
    dataIndex: 'msg',
    key: 'msg',
  },
]

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      msg: '',
      userName: '',
      loginModal: false,
      msglist: [],
    }
  }
  ws = null;
  getRandomInt = (max) => {
      return Math.floor(Math.random() * Math.floor(max));
  };
  getUrlVars = (parameter) => {
      var vars = {};
      window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m,key,value) {
          vars[key] = value;
      });
      return vars[parameter];
  };
  createWs = (userName) => {
      var roomId = this.getUrlVars('roomid');
      if (roomId === undefined) {
          roomId = new Date().getTime().toString() + "." + this.getRandomInt(100000).toString() + this.getRandomInt(100000).toString() + this.getRandomInt(100000).toString()
          window.history.pushState(null, null, window.location.href + '?roomid=' + roomId)
      }
      console.log('roomId:', roomId)
      this.ws = new WebSocket(`ws://127.0.0.1/chat?username=${userName}&roomid=${roomId}`);
      this.ws.onmessage = this.onMessage;
  };
  addMsg = (senderid, sendername, msg) => {
    var c = this.state.msglist.length
    var m = {'key': c, 'userid': senderid, 'username': sendername, 'msg': msg, 'datetime': this.nowStr()};
    this.setState({msglist:[...this.state.msglist, m]});
  };
  nowStr = () => {
    console.log(dayjs().format("MM-DD hh:mm:ss"))
    return dayjs().format('MM-DD hh:mm:ss');
  };
  onMessage = (evt) => {
    console.log("Received Message: " + evt.data);
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
  };
  handleLogin = () => {
    var userName = this.state.userName.trim();
    if (userName == '') {
      message.warning('用户名不能为空');
      return;
    }
    this.setState({loginModal: false})
    this.createWs(userName)
  };
  handleSend = () => {
    var value = this.state.msg.trim();
    if (value === "") {
      message.warning('消息内容为空，不能发送');
      return;
    }
    this.addMsg('', '我', value)
    this.ws.send(JSON.stringify({'type':'msg', 'msg': value}))
    this.setState({msg: ''})
    this.msgInput.focus();
  };
  handleEnter = (e) => {
    if (e.keyCode === 13 && e.ctrlKey) {
      this.handleSend()
      e.preventDefault()
      return
    }
  }
  componentWillMount() {
  }
  componentDidMount() {
    this.setState({loginModal: true})
  }
  render() {
    return (
      <div>
        <Modal
          title="输入用户名"
          visible={this.state.loginModal}
          closable={false}
          mask={true}
          maskClosable={false}
          footer={[
            <Button key="submit" type="primary" onClick={this.handleLogin}>登录</Button>
          ]}
        >
          <Input value={this.state.userName} onChange={(e)=>{this.setState({userName:e.target.value})}}></Input>
        </Modal>
        <Layout style={{minHeight: '100vh'}}>
          <Sider>Sider</Sider>
          <Layout>
            <Content>
              <Table columns={columns} dataSource={this.state.msglist} scroll={{y:400}} />
            </Content>
            <Footer style={{padding:0}}>
              <Layout>
                <div style={{display:'inline-flex',alignItems:'flex-end'}}>
                  <TextArea
                    style={{flex:'auto'}}
                    ref={(input) => { this.msgInput = input; }} 
                    onChange={(e)=>{this.setState({msg:e.target.value})}}
                    onKeyDown={this.handleEnter}
                    value={this.state.msg}
                    autoSize={{minRows:1, maxRows:4}}
                  />
                  <Button
                    style={{flex:'none'}}
                    type="primary"
                    onClick={this.handleSend}
                  >提交</Button>
                </div>
              </Layout>
            </Footer>
          </Layout>
        </Layout>

      </div>
    );
  }
}

export default App;
