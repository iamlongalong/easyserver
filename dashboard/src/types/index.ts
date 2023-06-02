
enum FileType {
    File = "file",
    Folder = "folder",
    Image = "image",
    Text = "text",
    Video = "video",
    Audio = "audio",
}

interface FileInfo {
    filename: string;
    filetype: FileType;
    modify_time: Date;
    size: string;
    fullpath: string;
}

interface ResFileInfo {
    name: string;
    size: number;
    is_dir: boolean;
    mod_time: number;
    file_type: FileType;
}


export {
    FileInfo,
    FileType,
    ResFileInfo
}
