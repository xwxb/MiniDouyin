package ffmpeg

import (
	"github.com/xwxb/MiniDouyin/config"
	"os/exec"
)

// Ffmpeg 传入视频地址和要输出的图片地址进行截图，方式是执行命令语句，前提要求安装ffmpeg并配置系统环境变量
// 需要在config的hide.go里补充一个自己本地项目的路径，如 const Path = "F:/workplace/MiniDouyin"
// 如果出现exec: file does not exist报错，请尝试执行 go env -w GOOS=windows，原因可能是之前配置过goos的linux环境变量
func Ffmpeg(videoName string, imageName string) (err error) {
	cmd := exec.Command("ffmpeg", "-ss", "00:00:01", "-i", config.Path + "/public/" + videoName, "-vframes", "1", config.Path + "/public/" + imageName)
	//var out bytes.Buffer
	//var stderr bytes.Buffer
	//cmd.Stdout = &out
	//cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		//fmt.Println(err.Error(), stderr.String())
		return err
	}
	//fmt.Println(out.String())
	return nil
}