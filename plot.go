package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"

	"gonum.org/v1/plot"

	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Plot the trajectory of the optimization.
// Require gonum/plot and support only 2-Dimentional.
// Return an error if something goes wrong, but it does not
// mean necessarily stop the whole program.
func plot2D(traj [IterMax][N][]float64) error {
	const dirPath string = "./traj/"
	if err := os.RemoveAll(dirPath); err != nil {
		return err
	}
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	var saveFiles = []string{}

	for iter := 0; iter < IterMax; iter++ {
		p := plot.New()
		p.Add(plotter.NewGrid())

		p.X.Label.Text = "x1"
		p.Y.Label.Text = "x2"
		p.X.Min = Xmin[0]
		p.X.Max = Xmax[0]
		p.Y.Min = Xmin[1]
		p.Y.Max = Xmax[1]

		pts := make(plotter.XYs, N)

		for i := 0; i < N; i++ {
			pts[i].X = traj[iter][i][0]
			pts[i].Y = traj[iter][i][1]
		}

		s, err := plotter.NewScatter(pts)
		if err != nil {
			return err
		}
		s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

		p.Add(s)

		// Save to a png file
		filename := fmt.Sprintf("%v/traj_%05d.png", dirPath, iter)
		if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
			return err
		}

		// Keep file names to generate GIF animation later
		saveFiles = append(saveFiles, filename)
	}

	// Generate GIF animation of the trajectory
	if err := encodeToGif(saveFiles); err != nil {
		return err
	}

	return nil
}

// Encode png files to gif animation
func encodeToGif(files []string) error {
	out := &gif.GIF{}
	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}

		imgSrc, _, err := image.Decode(f)
		if err != nil {
			return err
		}

		buf := bytes.Buffer{}

		if err = gif.Encode(&buf, imgSrc, nil); err != nil {
			return err
		}

		img, err := gif.Decode(&buf)
		if err != nil {
			return err
		}
		f.Close()

		out.Image = append(out.Image, img.(*image.Paletted))
		out.Delay = append(out.Delay, 0)
	}

	f, err := os.OpenFile("traj.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()
	gif.EncodeAll(f, out)
	return nil
}
