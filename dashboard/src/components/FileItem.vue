<template>
  <div>
    <div
      style="padding-top: 100%; position: relative"
      @dblclick="handleDBLClick"
    >
      <!-- 根据 filetype 改变图片地址 -->
      <el-image
        class="image-1-1"
        size="small"
        :src="getimage(file.filetype)"
        fit="contain"
        :lazy="true"
      ></el-image>
    </div>
    <div>
      <div class="file-name-info">
        <a class="notlink" :href="getFileDownloadLink(file.filename)">{{
          file.filename
        }}</a>
      </div>
      <div class="file-modify-info">
        {{ convertDate(file.modify_time) }}
      </div>
    </div>
  </div>
</template>


<script lang="ts">
import { Options, Vue } from "vue-class-component";
import { FileInfo, FileType } from "@/types";
import { getEnvUrl, getCurrDirPath, joinPath } from "@/utils";

let file: FileInfo;

@Options({
  props: {
    file: {
      required: true,
    },
  },
})
export default class FileItem extends Vue {
  file: FileInfo = file;

  // convert Date to human readable string
  convertDate(date: Date): string {
    return date.toLocaleDateString() + " " + date.toLocaleTimeString();
  }

  handleDBLClick() {
    this.$emit("fileDBClick", this.file);
  }

  getimage(t: FileType) {
    let src = `static/filepic_${t}.png`;

    if (t == FileType.Image) {
      src = getEnvUrl() + this.file.fullpath;
    }

    return src;
  }

  getFileDownloadLink(fname: string): string {
    let originUrl = getEnvUrl();
    let dirPath = getCurrDirPath();
    let link = joinPath(originUrl, dirPath, fname);

    return link;
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="less" scoped>
// 文字 不换行，超出后省略，居中显示
.file-name-info {
  font-size: 0.9rem;
  color: #333;
  white-space: nowrap;
  text-align: center;
  width: 95%;
  padding: 0 5px;
  margin: 5px 0;
  overflow: hidden;
  text-overflow: ellipsis;
}
.file-modify-info {
  font-size: 0.7rem;
  color: #999;
  margin-bottom: 5px;
}
.image-1-1 {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.notlink {
  text-decoration: none;
  color: inherit;

  &:hover {
    text-decoration: underline;
  }
}
</style>
