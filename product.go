package docweaver

import "fmt"

type Product struct {
	Name        string
	Description string
	ImageUrl    string
	Versions    []string
	Directory   string
}

type productRoot struct {
	ParentDir string
	Name      string
	Source    string
}

func (p *productRoot) fullPath() string {
	return fmt.Sprintf("%s/%s", p.ParentDir, p.Name)
}
