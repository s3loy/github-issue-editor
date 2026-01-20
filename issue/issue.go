package issue

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Issue struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	State string `json:"state,omitempty"`
}

func getAuthToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("未读取到GITHUB_TOKEN 环境变量")
	}
	return token
}

func getEditorInput(init string) (string, error) {
	tmpFile, err := os.CreateTemp("", "issue_*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if init != "" {
		if _, err = tmpFile.WriteString(init); err != nil {
			return "", err
		}
		tmpFile.Sync()
	}
	tmpFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "Code"
	}

	fmt.Printf("正在启动编辑器 (%s)...\n", editor)

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return "", fmt.Errorf("启动编辑器失败: %v", err)
	}
	fmt.Println("------------------------------------------------------")
	fmt.Println("请在编辑器中修改内容，完成后保存并关闭文件。")
	fmt.Print(">>> 确认提交请输入 1，输入其他内容将取消操作: ")

	scanner := bufio.NewScanner(os.Stdin)
	var confirm string
	if scanner.Scan() {
		confirm = strings.TrimSpace(scanner.Text())
	}

	if confirm != "1" {
		return "", fmt.Errorf("用户取消操作")
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(content)), nil
}

func CreateIssue(owner, repo string) error {
	fmt.Print("请输入issue标题:")
	scanner := bufio.NewScanner(os.Stdin)
	var title string
	if scanner.Scan() {
		title = scanner.Text()
	}

	if title == "" {
		return fmt.Errorf("标题不能为空")
	}

	body, err := getEditorInput("")
	if err != nil {
		return err
	}

	issue := Issue{Title: title, Body: body}
	data, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "token "+getAuthToken())
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		bodyRecieved, _ := io.ReadAll(response.Body)
		return fmt.Errorf("创建issue失败: 返回为 \n%s", string(bodyRecieved))
	}

	var result map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&result)
	fmt.Printf("Issue 创建成功! 编号: #%.0f\n", result["number"])

	return nil
}

func ReadIssue(owner, repo string, number int) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", owner, repo, number)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "token "+getAuthToken())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	bodyRecieved, err := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("读取issue失败: 返回为 \n%s", string(bodyRecieved))
	}
	if err != nil {
		return err
	}
	fmt.Println("Issue信息: ")
	fmt.Println(string(bodyRecieved))
	return nil
}

func UpdateIssue(owner, repo string, number int) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", owner, repo, number)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "token "+getAuthToken())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyRecieved, _ := io.ReadAll(response.Body)
		return fmt.Errorf("读取issue失败: %s", string(bodyRecieved))
	}

	var current Issue
	if err = json.NewDecoder(response.Body).Decode(&current); err != nil {
		return err
	}

	fmt.Printf("当前标题为: %s\n", current.Title)
	fmt.Print("请输入新的标题(留空不变): ")

	scanner := bufio.NewScanner(os.Stdin)
	var newTitle string
	if scanner.Scan() {
		newTitle = scanner.Text()
	}
	if newTitle == "" {
		newTitle = current.Title
	}

	newBody, err := getEditorInput(current.Body)
	if err != nil {
		return err
	}

	updatedIssue := Issue{Title: newTitle, Body: newBody}
	data, err := json.Marshal(updatedIssue)
	if err != nil {
		return err
	}

	request, err = http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "token "+getAuthToken())
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		bodyRecieved, _ := io.ReadAll(response.Body)
		return fmt.Errorf("更新issue失败: 返回为\n%s", string(bodyRecieved))
	}
	fmt.Println("Issue 更新成功")
	return nil
}

func CloseIssue(owner, repo string, number int) error {
	updatedIssue := Issue{State: "closed"}
	data, err := json.Marshal(updatedIssue)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", owner, repo, number)
	request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "token "+getAuthToken())
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		bodyRecieved, _ := io.ReadAll(response.Body)
		return fmt.Errorf("关闭issue失败: 返回为\n%s", string(bodyRecieved))
	}
	fmt.Println("Issue已关闭")
	return nil
}
