package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type OSM struct {
	XMLName xml.Name `xml:"osm"`
	Ways    []Way    `xml:"way"`
	Nodes   []Node   `xml:"node"`
}

type Way struct {
	ID   uint64 `xml:"id,attr"`
	NDs  []ND   `xml:"nd"`
	Tags []Tag  `xml:"tag"`
}

type ND struct {
	Ref uint64 `xml:"ref,attr"`
}

type Node struct {
	ID  uint64  `xml:"id,attr"`
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
}

type Tag struct {
	K string `xml:"k,attr"`
	V string `xml:"v,attr"`
}

// BBox represents a geographical bounding box
type BBox struct {
	MinLat float64
	MinLon float64
	MaxLat float64
	MaxLon float64
}

func main() {
	// print the current dir
	dir, err := os.Getwd()
	fmt.Println("dir:", dir)
	osmFile := "./resources/200.osm" // 更改为你的 OSM XML 文件名
	data, err := ioutil.ReadFile(osmFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var osm OSM
	err = xml.Unmarshal(data, &osm)
	if err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return
	}

	// 创建节点 ID 到经纬度的映射
	nodeMap := make(map[uint64]Node)
	for _, node := range osm.Nodes {
		nodeMap[node.ID] = node
	}

	bbox := BBox{
		MinLat: 180.0,
		MinLon: 180.0,
		MaxLat: -180.0,
		MaxLon: -180.0,
	}

	// 遍历所有 way 元素
	for _, way := range osm.Ways {
		for _, tag := range way.Tags {
			// 检查当前 way 是否含有特定标签
			if tag.K == "LINE_TYPE" && tag.V == "RoadCenter" {
				// 如果含有特定标签，则计算该 way 的边界框
				for _, nd := range way.NDs {
					node := nodeMap[nd.Ref]
					if node.Lat < bbox.MinLat {
						bbox.MinLat = node.Lat
					}
					if node.Lon < bbox.MinLon {
						bbox.MinLon = node.Lon
					}
					if node.Lat > bbox.MaxLat {
						bbox.MaxLat = node.Lat
					}
					if node.Lon > bbox.MaxLon {
						bbox.MaxLon = node.Lon
					}
				}
				break // 不需要检查其他标签
			}
		}
	}

	// 打印边界框结果
	fmt.Printf("BBox: %+v\n", bbox)
}
