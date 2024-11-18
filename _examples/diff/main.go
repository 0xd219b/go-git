package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// createTestRepo 创建一个临时的测试仓库
func createTestRepo() (string, error) {
	// 创建临时目录
	dir, err := ioutil.TempDir("", "git-test-*")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}

	// 初始化仓库
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		os.RemoveAll(dir)
		return "", fmt.Errorf("init repo: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		os.RemoveAll(dir)
		return "", fmt.Errorf("get worktree: %w", err)
	}

	// 创建测试文件
	testFiles := map[string]string{
		"file1.txt": "Original content\nfor file1\n",
		"file2.txt": "Original content\nfor file2\n",
	}

	for name, content := range testFiles {
		path := filepath.Join(dir, name)
		if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
			os.RemoveAll(dir)
			return "", fmt.Errorf("write file: %w", err)
		}

		_, err = w.Add(name)
		if err != nil {
			os.RemoveAll(dir)
			return "", fmt.Errorf("git add: %w", err)
		}
	}

	// 创建初始提交
	_, err = w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Author",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		os.RemoveAll(dir)
		return "", fmt.Errorf("create commit: %w", err)
	}

	// 修改文件以产生差异
	modifiedFiles := map[string]string{
		"file1.txt": "Modified content\nfor file1\nwith new line\n",
		"file3.txt": "New file content\n",
	}

	for name, content := range modifiedFiles {
		path := filepath.Join(dir, name)
		if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
			os.RemoveAll(dir)
			return "", fmt.Errorf("write modified file: %w", err)
		}
	}

	// 删除 file2.txt 以测试删除场景
	if err := os.Remove(filepath.Join(dir, "file2.txt")); err != nil {
		os.RemoveAll(dir)
		return "", fmt.Errorf("remove file: %w", err)
	}

	return dir, nil
}

func main() {
	// 创建测试仓库
	dir, err := createTestRepo()
	if err != nil {
		fmt.Printf("创建测试仓库失败: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir) // 清理临时目录

	// 打开仓库
	repo, err := git.PlainOpen(dir)
	if err != nil {
		fmt.Printf("无法打开仓库: %v\n", err)
		os.Exit(1)
	}

	// 获取工作区
	w, err := repo.Worktree()
	if err != nil {
		fmt.Printf("无法获取工作区: %v\n", err)
		os.Exit(1)
	}

	// 打印差异
	err = w.PrintDiff(&git.DiffOptions{
		DetectBinary: true,
	})
	if err != nil {
		fmt.Printf("打印差异失败: %v\n", err)
		os.Exit(1)
	}
}
