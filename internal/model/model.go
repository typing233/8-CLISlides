package model

type Metadata struct {
	Title    string `yaml:"title"`
	Author   string `yaml:"author"`
	Date     string `yaml:"date"`
	Theme    string `yaml:"theme"`
	Pager    bool   `yaml:"pager"`
	SSHPort  int    `yaml:"ssh_port"`
	SSHHost  string `yaml:"ssh_host"`
}

type Slide struct {
	Raw      string
	Rendered string
	Index    int
}

type Presentation struct {
	Meta   Metadata
	Slides []Slide
}
