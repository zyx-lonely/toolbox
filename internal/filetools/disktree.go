package filetools

import (
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// DiskTreeNode 磁盘空间节点
type DiskTreeNode struct {
	Name     string         `json:"name"`
	Path     string         `json:"path"`
	Size     uint64         `json:"size"`
	IsDir    bool           `json:"isDir"`
	Children []DiskTreeNode `json:"children,omitempty"`
}

// ScanDiskTree 扫描目录并返回树状结构
func ScanDiskTree(rootPath string, maxDepth int) (*DiskTreeNode, error) {
	if maxDepth <= 0 {
		maxDepth = 3
	}

	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}

	node := &DiskTreeNode{
		Name:  info.Name(),
		Path:  rootPath,
		Size:  0,
		IsDir: info.IsDir(),
	}

	if info.IsDir() {
		node.Children = scanDir(rootPath, maxDepth, 1)
		// 计算总大小
		for _, child := range node.Children {
			node.Size += child.Size
		}
	}

	return node, nil
}

func scanDir(dirPath string, maxDepth, currentDepth int) []DiskTreeNode {
	if currentDepth > maxDepth {
		return nil
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil
	}

	var nodes []DiskTreeNode
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 限制并发数
	sem := make(chan struct{}, 20)

	for _, entry := range entries {
		// 跳过隐藏文件和系统文件
		if entry.Name()[0] == '.' || entry.Name() == "System Volume Information" {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		node := DiskTreeNode{
			Name:  entry.Name(),
			Path:  fullPath,
			Size:  uint64(info.Size()),
			IsDir: entry.IsDir(),
		}

		if entry.IsDir() {
			wg.Add(1)
			sem <- struct{}{}
			go func(n DiskTreeNode, p string) {
				defer wg.Done()
				defer func() { <-sem }()
				n.Children = scanDir(p, maxDepth, currentDepth+1)
				// 计算子目录总大小
				for _, child := range n.Children {
					n.Size += child.Size
				}
				mu.Lock()
				nodes = append(nodes, n)
				mu.Unlock()
			}(node, fullPath)
		} else {
			mu.Lock()
			nodes = append(nodes, node)
			mu.Unlock()
		}
	}

	wg.Wait()

	// 按大小降序排序
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Size > nodes[j].Size
	})

	// 限制返回数量，避免数据过大
	if len(nodes) > 50 {
		nodes = nodes[:50]
	}

	return nodes
}
