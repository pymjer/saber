package graph

import "container/list"

type Edge struct {
	sid int // 边的开始顶点
	tid int // 边的终止顶点
	w   int // 权重
}

type Vertex struct {
	id   int    // 顶点编号ID
	name string // 顶点名称
}

type Graph struct {
	adj []*list.List // 领接表
	v   int          // 顶点个数
}

func (g *Graph) AddEdge()  {
	
}
