package rt

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type ObjFile struct {
	Ignored			int
	Vertices 		[]*Tuple
	Normals 		[]*Tuple
	DefaultGroup 	*Shape
	Groups			map[string]*Shape
}


func ParseObjFile(path string) (*ObjFile, error) {
	file, err := os.Open(path)
	if err != nil {return nil, err}
	defer file.Close()

	return ParseObjFileFromScanner(bufio.NewScanner(file))
}

func ParseObjFileFromString(content string) (*ObjFile, error) {
	return ParseObjFileFromScanner(bufio.NewScanner(strings.NewReader(content)))
}

func ParseObjFileFromScanner(scanner *bufio.Scanner) (*ObjFile, error) {
	objFile := &ObjFile{
		Ignored: 0,
		Vertices: make([]*Tuple, 0),
		DefaultGroup: NewGroup(),
		Groups: make(map[string]*Shape),
	}

	currentGroup := objFile.DefaultGroup

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			objFile.Ignored++
			continue
		}

		fields := strings.Fields(line)

		switch fields[0] {
		case "v":
			objFile.Vertices = append(objFile.Vertices, parseVertex(fields[1:]))
		case "vn":
			objFile.Normals = append(objFile.Normals, parseNormal(fields[1:]))
		case "f":
			currentGroup.AddChildren(parseFace (objFile, fields[1:])...)
		case "g":
			currentGroup = NewGroup()
			objFile.Groups[fields[1]] = currentGroup
		default:
			objFile.Ignored++
		}
	}
	return objFile, nil
}

func parseFloat (s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		f = 0
	}
	return f
}

func parseInt (s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		i = 0
	}
	return i
}

func parseTriple (fields []string, fn func(x, y, z float64) *Tuple) *Tuple {
	if len(fields) < 3 {
		return nil
	}
	x := parseFloat(fields[0])
	y := parseFloat(fields[1])
	z := parseFloat(fields[2])
	return fn(x,y,z)
}

func parseVertex(fields []string) *Tuple {
	return parseTriple(fields, NewPoint)
}

func parseNormal (fields []string) *Tuple {
	return parseTriple(fields, NewVector)
}

func parseFace (objFile *ObjFile, fields []string) []*Shape {
	if len(fields) < 3{
		return nil
	}
	verIdxs := parseVertexIndices (fields)
	normIdxs := parseNormalIndices (fields)

	if len (normIdxs) == len (verIdxs) {
		return buildSmoothTriangles (objFile, verIdxs, normIdxs)
	}
	return buildSimpleTriangles(objFile, verIdxs)
}

func parseVertexIndices (fields []string) []int {
	vals := make([]int, 0)
	for _, f := range fields {
		parts := strings.Split(f, "/")
		vals = append(vals, parseInt(parts[0])-1)
	}
	return vals
}


func parseNormalIndices (fields []string) []int {
	vals := make([]int, 0)

	for _, f := range fields {
		parts := strings.Split(f, "/")
		if len(parts) >= 3 {
			vals = append(vals, parseInt(parts[2])-1)
		}
	}
	return vals
}

func buildSimpleTriangles (objFile *ObjFile, verIdxs []int) []*Shape{
	log.Printf("building a simple triangle")
	shapes := make ([]*Shape, 0)
	a := verIdxs[0]
	for i := 1; i < len(verIdxs) - 1; i++ {
		b := verIdxs[i]
		c := verIdxs[i+1]

		shapes = append(shapes, NewTriangle(objFile.Vertices[a], objFile.Vertices[b], objFile.Vertices[c]))
	}
	return shapes
}

func buildSmoothTriangles (objFile *ObjFile, verIdxs, normIdxs []int) []*Shape{
	log.Printf("building a smooth triangle")
	shapes := make ([]*Shape, 0)
	a := verIdxs[0]
	x := normIdxs[0]
	for i := 1; i < len(verIdxs) - 1; i++ {
		b := verIdxs[i]
		c := verIdxs[i+1]
		y := normIdxs[i]
		z := normIdxs[i+1]

		shapes = append(shapes, NewSmoothTriangle(
			objFile.Vertices[a], objFile.Vertices[b], objFile.Vertices[c],
			objFile.Normals[x], objFile.Normals[y], objFile.Normals[z]))
	}
	return shapes
}

func (o *ObjFile) ToGroup () *Shape {
	g := NewGroup()
	g.AddChildren(o.DefaultGroup)
	for _, c := range o.Groups {
		g.AddChildren(c)
	}
	return g
}