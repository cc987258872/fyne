package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectEntry_MinSize(t *testing.T) {
	smallOptions := []string{"A", "B", "C"}

	largeOptions := []string{"Large Option A", "Larger Option B", "Very Large Option C"}
	largeOptionsMinWidth := optionsMinSize(largeOptions).Width

	minTextHeight := widget.NewLabel("W").MinSize().Height

	tests := map[string]struct {
		placeholder string
		value       string
		options     []string
		want        fyne.Size
	}{
		"empty": {
			want: fyne.NewSize(emptyTextWidth()+4*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"empty + small options": {
			options: smallOptions,
			want:    fyne.NewSize(emptyTextWidth()+dropDownIconWidth()+4*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"empty + large options": {
			options: largeOptions,
			want:    fyne.NewSize(largeOptionsMinWidth+2*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"value": {
			value: "foo",
			want:  widget.NewLabel("foo").MinSize().Add(fyne.NewSize(4*theme.Padding(), 2*theme.Padding())),
		},
		"large value + small options": {
			value:   "large",
			options: smallOptions,
			want:    widget.NewLabel("large").MinSize().Add(fyne.NewSize(dropDownIconWidth()+4*theme.Padding(), 2*theme.Padding())),
		},
		"small value + large options": {
			value:   "small",
			options: largeOptions,
			want:    fyne.NewSize(largeOptionsMinWidth+2*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"placeholder": {
			placeholder: "example",
			want:        widget.NewLabel("example").MinSize().Add(fyne.NewSize(4*theme.Padding(), 2*theme.Padding())),
		},
		"large placeholder + small options": {
			placeholder: "large",
			options:     smallOptions,
			want:        widget.NewLabel("large").MinSize().Add(fyne.NewSize(dropDownIconWidth()+4*theme.Padding(), 2*theme.Padding())),
		},
		"small placeholder + large options": {
			placeholder: "small",
			options:     largeOptions,
			want:        fyne.NewSize(largeOptionsMinWidth+2*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := widget.NewSelectEntry(tt.options)
			e.PlaceHolder = tt.placeholder
			e.Text = tt.value
			assert.Equal(t, tt.want, e.MinSize())
		})
	}
}

func TestSelectEntry_DropDown(t *testing.T) {
	options := []string{"A", "B", "C"}
	e := widget.NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	c := w.Canvas()

	assert.Nil(t, c.Overlays().Top())

	var dropDownSwitch *widget.Button
	for _, o := range test.LaidOutObjects(c.Content()) {
		if b, ok := o.(*widget.Button); ok {
			dropDownSwitch = b
			break
		}
	}
	require.NotNil(t, dropDownSwitch, "drop down switch not found")

	test.Tap(dropDownSwitch)
	require.NotNil(t, c.Overlays().Top(), "drop down didn't open")
	require.IsType(t, &widget.PopUp{}, c.Overlays().Top(), "drop down is not a *widget.PopUp")

	popUp := c.Overlays().Top().(*widget.PopUp)
	entryMinWidth := dropDownIconWidth() + emptyTextWidth() + 4*theme.Padding()
	assert.Equal(t, optionsMinSize(options).Max(fyne.NewSize(entryMinWidth-2*theme.Padding(), 0)), popUp.Content.Size())
	assert.Equal(t, options, popUpOptions(popUp), "drop down menu texts don't match SelectEntry options")

	tapPopUpItem(t, popUp, 1)
	assert.Nil(t, c.Overlays().Top())
	assert.Equal(t, "B", e.Text)

	test.Tap(dropDownSwitch)
	popUp = c.Overlays().Top().(*widget.PopUp)
	tapPopUpItem(t, popUp, 2)
	assert.Nil(t, c.Overlays().Top())
	assert.Equal(t, "C", e.Text)
}

func dropDownIconWidth() int {
	dropDownIconWidth := theme.IconInlineSize() + theme.Padding()
	return dropDownIconWidth
}

func emptyTextWidth() int {
	return widget.NewLabel("M").MinSize().Width
}

func optionsMinSize(options []string) fyne.Size {
	var labels []*widget.Label
	for _, option := range options {
		labels = append(labels, widget.NewLabel(option))
	}
	minWidth := 0
	minHeight := 0
	for _, label := range labels {
		if minWidth < label.MinSize().Width {
			minWidth = label.MinSize().Width
		}
		minHeight += label.MinSize().Height
	}
	// padding between all options
	minHeight += (len(labels) - 1) * theme.Padding()
	return fyne.NewSize(minWidth, minHeight)
}

func popUpOptions(popUp *widget.PopUp) []string {
	var texts []string
	for _, o := range test.LaidOutObjects(popUp.Content) {
		if t, ok := o.(*canvas.Text); ok {
			texts = append(texts, t.Text)
		}
	}
	return texts
}

func tapPopUpItem(t *testing.T, p *widget.PopUp, i int) {
	var items []fyne.Tappable
	for _, o := range test.LaidOutObjects(p.Content) {
		if t, ok := o.(fyne.Tappable); ok {
			items = append(items, t)
		}
	}
	require.Greater(t, len(items), i, "not enough tappables found (%d out of at least %d)", len(items), i+1)
	test.Tap(items[i])
}
