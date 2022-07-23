package ui

type Bound struct {
	begin int
	end   int
}

type Container struct {
	xBound  Bound
	yBound  Bound
	xCursor int
	yCursor int
	Size    int
	isFull  bool
}

func NewContainer(x, y int) *Container {
	x_Bound := Bound{
		begin: 20,
		end:   x - (x / 4),
	}
	y_Bound := Bound{
		begin: 10,
		end:   y - 5,
	}
	Size := (x_Bound.end - x_Bound.begin) * (y_Bound.end - y_Bound.begin)
	return &Container{
		xBound:  x_Bound,
		yBound:  y_Bound,
		xCursor: x_Bound.begin,
		yCursor: y_Bound.begin,
		Size:    Size,
		isFull:  false,
	}

}
