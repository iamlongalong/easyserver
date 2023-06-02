<template>
  <el-container height="100%" direction="vertical">
    <!-- <el-header height="100%"> header </el-header> -->
    <el-main height="100%">
      <!-- 顶部面包屑导航 -->
      <el-row :gutter="20" justify="space-around">
        <el-col :span="12" :offset="0">
          <!-- 导航栏 -->
          <!-- <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">/</el-breadcrumb-item>
            <el-breadcrumb-item><a href="/">data</a></el-breadcrumb-item>
            <el-breadcrumb-item>xxxx</el-breadcrumb-item>
          </el-breadcrumb> -->
        </el-col>
        <el-col :span="12" :offset="0">
          <!-- 上传按钮 -->
          <el-button type="primary" size="large" @click="upload.show = true"
            ><el-icon><UploadFilled /></el-icon>上传</el-button
          >
        </el-col>
      </el-row>

      <!-- 文件展示块 -->
      <el-row :gutter="20" justify="center">
        <el-col :span="16" :offset="0">
          <el-row justify="center">
            <el-col :span="12" :offset="0">
              <h2>
                <span v-for="(path, index) in currentPaths" :key="index">
                  <span @click="goToPathByIndex(index)" class="can-click"
                    >{{ path }}
                  </span>
                  <span v-if="path != '/'">/</span>
                </span>
              </h2>
            </el-col>
          </el-row>

          <div
            style="height: 60vh; position: relative; background-color: #efefef"
          >
            <div style="margin: 15px 0 15px 15px; height: inherit">
              <div v-if="files.length == 0" style="height: inherit">
                <el-row
                  :gutter="20"
                  justify="center"
                  align="middle"
                  style="height: inherit"
                >
                <el-col :offset="0">
                  <div style="padding: 20px 0;">
                    <p style="font-size: 1.5rem; color: #aaa;">no files in this folder</p>
                  </div>
  
                  <div>
                    <el-button type="primary" @click="goBack"
                      >返回上一级</el-button
                    >
                  </div>
                </el-col>
                
                </el-row>
              </div>

              <div
                v-if="files.length != 0"
                style="height: 60vh; overflow-y: auto; overflow-x: hidden"
              >
                <el-row :gutter="20" justify="start">
                  <div
                    v-for="(file, index) in files"
                    :key="index"
                    style="margin: 5px 10px"
                  >
                    <!-- files block -->
                    <FileItem
                      @fileDBClick="handledblClick"
                      style="width: 140px"
                      :file="file"
                    ></FileItem>
                  </div>
                </el-row>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>

      <!-- 上传展示块 -->
      <el-dialog
        @close="handleClose"
        v-model="upload.show"
        title=""
        width="80%"
        align-center
      >
        <div>
          <el-upload
            ref="uploadref"
            :on-success="handleSuccess"
            show-file-list
            drag
            multiple
            :http-request="uploadFile"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              Drop file here or <em>click to upload</em>
            </div>
          </el-upload>
        </div>
      </el-dialog>

      <!-- 底部的认证方式 -->

      <el-affix
        @click="auth.show = !auth.show"
        class="affix"
        position="bottom"
        :offset="50"
      >
        <div style="width: 30px">
          <el-avatar
            icon="el-icon-user-solid"
            size="large"
            shape="circle"
            :src="userAvatar"
            fit="contain"
          ></el-avatar>
        </div>
      </el-affix>

      <el-dialog v-model="auth.show" align-center style="max-width: 600px">
        <template #header>
          <el-menu
            @select="(index) => (auth.selected = index)"
            mode="horizontal"
            :default-active="auth.selected"
          >
            <el-menu-item index="basic">账号密码认证</el-menu-item>
            <el-menu-item index="token">TOKEN 认证</el-menu-item>
            <el-menu-item index="annymouns">匿名使用</el-menu-item>
          </el-menu>
        </template>
        <!-- <AuthBox></AuthBox> -->

        <!-- 认证内容 -->
        <div style="margin: 0 20px; height: 100px">
          <!-- 账号密码认证 -->
          <div v-if="auth.selected == 'basic'">
            <el-row :gutter="20" justify="center" align="middle">
              <el-col :span="18" :offset="0">
                <el-input
                  v-model="auth.basic.username"
                  placeholder="请输入用户名"
                >
                  <template #prepend>
                    <div style="width: 60px">username</div>
                  </template>
                </el-input>
              </el-col>
            </el-row>
            <br />
            <el-row :gutter="20" justify="center">
              <el-col :span="18" :offset="0">
                <el-input
                  label="password"
                  show-password
                  v-model="auth.basic.password"
                  placeholder="请输入密码"
                >
                  <template #prepend>
                    <div style="width: 60px">password</div>
                  </template>
                </el-input>
              </el-col>
            </el-row>
          </div>

          <!-- token 认证 -->
          <div v-if="auth.selected == 'token'">
            <el-row
              :gutter="20"
              justify="center"
              style="height: 100px"
              align="middle"
            >
              <el-col :span="18" :offset="0">
                <el-input v-model="auth.token" placeholder="请输入 token">
                  <template #prepend>
                    <div style="width: 60px">token</div>
                  </template>
                </el-input>
              </el-col>
            </el-row>
          </div>

          <!-- 匿名访问 -->
          <div v-if="auth.selected == 'annymouns'">
            <el-row
              :gutter="20"
              justify="center"
              style="height: 100px"
              align="middle"
              >使用匿名访问</el-row
            >
          </div>
        </div>

        <template #footer>
          <el-row :gutter="20" justify="center">
            <el-button type="primary" size="large" @click="confirmAuth">
              确定
            </el-button>
          </el-row>
        </template>
      </el-dialog>
    </el-main>
    <!-- <el-footer height="100px"> footer </el-footer> -->
  </el-container>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import AuthBox from "./AuthBox.vue";
import FileItem from "./FileItem.vue";
import UploadBox from "./UploadBox.vue";
import { FileInfo, FileType, ResFileInfo } from "@/types";
import { getEnvUrl } from "@/utils";
import {
  UploadFile,
  UploadFiles,
  UploadRequestOptions,
  UploadUserFile,
  UploadInstance,
  ElNotification,
  ElUpload,
  ElAlert,
  ElMessage,
} from "element-plus";
import { ref, h } from "vue";

import axios from "axios";

@Options({
  props: {},
  components: {
    AuthBox,
    FileItem,
    UploadBox,
  },
})
export default class HomePage extends Vue {
  auth = {
    show: false,
    basic: {
      username: "",
      password: "",
    },
    token: "",

    selected: "basic",
    useAuthType: "annymous",
  };

  upload = {
    show: false,
  };

  ischange = false;

  userAvatar = "/_dash/static/user.png";

  clearFiles() {
    const uploadref = this.$refs.uploadref as UploadInstance;

    uploadref.clearFiles();
  }

  handleClose() {
    setTimeout(this.clearFiles, 1000);

    if (this.ischange) {
      this.getPathFiles(this.currentDir).finally(() => {
        this.ischange = false;
      });
    }
  }

  currentPaths: string[] = ["/"];

  handledblClick(file: FileInfo) {
    if (file.filetype === FileType.Folder) {
      this.currentPaths.push(file.filename);
      this.getPathFiles(this.currentDir);
    } else {
      // open file
      window.open(getEnvUrl() + this.cleanPaths("/" + file.fullpath));
    }
  }

  confirmAuth() {
    if (this.auth.selected == "basic") {
      // 校验
      if (this.auth.basic.username == "" || this.auth.basic.password == "") {
        ElMessage({
          message: "用户名或密码不能为空",
          type: "error",
        });
        return;
      }
    } else if (this.auth.selected == "token") {
      // 校验
      if (this.auth.token == "") {
        ElMessage({
          message: "token 不能为空",
          type: "error",
        });
        return;
      }
    } else if (this.auth.selected == "annymous") {
    }

    this.auth.show = false;
    this.auth.useAuthType = this.auth.selected;

    // 存储到 localstorage
    localStorage.setItem("auth", JSON.stringify(this.auth));

    this.getPathFiles(this.currentDir);
  }

  goBack() {
    this.currentPaths.pop();
    this.getPathFiles(this.currentDir);
  }

  handleSuccess(
    response: any,
    uploadFile: UploadFile,
    uploadFiles: UploadFiles
  ) {
    // 发提示消息
    ElNotification({
      // title: "",
      duration: 3000,
      message: h(
        "i",
        { style: "color: teal" },
        `文件 ${uploadFile.name} 上传成功`
      ),
    });
  }

  goToPathByIndex(index: number) {
    this.currentPaths = this.currentPaths.slice(0, index + 1);
    this.getPathFiles(this.currentDir);
  }

  cleanPaths(path: string): string {
    return path.replace(/\/+/g, "/");
  }

  get currentDir(): string {
    let rawPath = this.currentPaths.join("/");
    // clean rawPath
    return rawPath.replace(/\/+/g, "/");
  }

  created() {
    // 从 localstorage 中获取 auth
    let authStr = localStorage.getItem("auth");
    // 如果有 auth, 则解析到 this.auth
    if (authStr) {
      this.auth = JSON.parse(authStr);
      this.getPathFiles(this.currentDir);
    } else {
      // 否则弹出 auth 对话框
      this.auth.show = true;
    }
  }

  // convert byte nums to human readable string
  private convertBytes(bytes: number): string {
    if (bytes < 1024) {
      return bytes + "B";
    } else if (bytes < 1024 * 1024) {
      return (bytes / 1024).toFixed(2) + "KB";
    } else if (bytes < 1024 * 1024 * 1024) {
      return (bytes / 1024 / 1024).toFixed(2) + "MB";
    } else {
      return (bytes / 1024 / 1024 / 1024).toFixed(2) + "GB";
    }
  }

  // use axios to get FileInfos from server
  getPathFiles(path: string): Promise<any> {
    let baseUrl = "/_apilist";

    // join url path with baseUrl and currentDir, clean then
    let finalpath = [baseUrl, path].join("/").replace(/\/+/g, "/");

    finalpath = getEnvUrl() + finalpath;

    // 使用认证
    if (this.auth.useAuthType == "basic") {
      let auth = btoa(
        this.auth.basic.username + ":" + this.auth.basic.password
      );
      axios.defaults.headers.common["Authorization"] = "Basic " + auth;
    } else if (this.auth.useAuthType == "token") {
      finalpath += "?token=" + this.auth.token;
    } else if (this.auth.useAuthType == "annymous") {
      axios.defaults.headers.common["Authorization"] = "";
    }

    return axios
      .get(finalpath)
      .then((res) => {
        let resFiles = res.data.data as ResFileInfo[];

        let newfiles: FileInfo[] = [];

        resFiles.forEach((resFile) => {
          newfiles.push({
            filename: resFile.name,
            filetype: resFile.file_type,
            modify_time: new Date(resFile.mod_time),
            size: this.convertBytes(resFile.size),
            fullpath: this.currentDir + "/" + resFile.name,
          });
        });

        this.files = newfiles;
      })
      .catch((err) => {
        ElMessage({
          showClose: true,
          message: "get files list error: " + err,
          type: "error",
        });
      });
  }

  changePath(path: string) {
    this.currentPaths.push(path);
  }

  uploadFile(options: UploadRequestOptions): Promise<XMLHttpRequest> {
    this.ischange = true;

    return new Promise((resolve, reject) => {
      let file = options.file;
      const xhr = new XMLHttpRequest();

      let url = getEnvUrl() + this.currentDir + "/" + file.name;

      debugger;
      // 使用认证
      if (this.auth.useAuthType == "basic") {
        let auth = btoa(
          this.auth.basic.username + ":" + this.auth.basic.password
        );
        xhr.open("POST", url, true);
        xhr.setRequestHeader("Authorization", "Basic " + auth);
      } else if (this.auth.useAuthType == "token") {
        url = url + "?token=" + this.auth.token;

        xhr.open("POST", url, true);
      } else {
        xhr.open("POST", url, true);
      }

      // Set up the request payload
      const formData = new FormData();
      formData.append("file", file);
      xhr.send(formData);

      // Handle the response
      xhr.onreadystatechange = () => {
        if (xhr.readyState === 4) {
          if (xhr.status === 200) {
            resolve(xhr);
          } else {
            // 提示出错
            ElMessage({
              showClose: true,
              dangerouslyUseHTMLString: true,
              message: `<p>upload file error:</p><br/>filename: ${file.name}<p>status: ${xhr.status}</p><p>statusText: ${xhr.statusText}</p>`,
              type: "error",
            });

            reject(xhr);
          }
        }
      };
    });
  }

  files: FileInfo[] = [
    {
      filename: "test",
      filetype: FileType.Folder,
      modify_time: new Date(),
      size: "4K",
    } as FileInfo,
    {
      filename: "test_file.txt",
      filetype: FileType.File,
      modify_time: new Date(),
      size: "10K",
    } as FileInfo,
    {
      filename: "test.jpg",
      filetype: FileType.Image,
      modify_time: new Date(),
      size: "1MB",
    } as FileInfo,
  ];
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="less" scoped>
div {
  // border: 0.1px solid #ebeef5;
}

h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}

.can-click {
  cursor: pointer;
  &:hover {
    color: #59b8fc;
  }
}

/* 垂直滚动条 */
::-webkit-scrollbar {
  width: 4px;
  height: 8px;
}

/* 滑块 */
::-webkit-scrollbar-thumb {
  background-color: #ccc;
  border-radius: 10px;
}

/* 滑道 */
::-webkit-scrollbar-track {
  background-color: #f1f1f1;
  border-radius: 10px;
}

.affix {
  position: fixed;
  bottom: 100px;
  right: 100px;

  &:hover {
    cursor: pointer;
  }
}
</style>
