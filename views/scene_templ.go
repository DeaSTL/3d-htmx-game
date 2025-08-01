// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.898
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/deastl/htmx-doom/gameobjects"
)

func FloorTile(x int, y int, width int, height int) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)

		style := fmt.Sprintf(`
      transform-origin: 0 0 0; 
      position: absolute; 
      transform: translate3d(%dpx,0px,%dpx) rotate3d(1,0,0,90deg); 
      background-image: url('/public/stone-floor.jpg'); 
      width: %dpx; 
      height: %dpx; 
      background-size: 255px 255px; 
      transform-style: preserve-3d;">
      `, x, y, width, height)
		templ_7745c5c3_Err = templ.Raw(`<div style="`+style+`"></div>`).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func Scene(game *gameobjects.GameMap) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		width := 16384
		height := 16384
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div id=\"scene\" class=\"scene\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, wall := range game.Walls {
			templ_7745c5c3_Err = Plane(wall).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		for x := 0; x < 2; x++ {
			for y := 0; y < 2; y++ {
				templ_7745c5c3_Err = FloorTile(x*width, y*height, width, height).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<div style=\"position: absolute; transform: scale3d(100,100,100) translate3d(0,1px,0) rotate3d(1,0,0,90deg); background-image: radial-gradient(#00BFFF, white, blue); width: 2000px; height: 2000em; transform-style: preserve-3d;\"></div><style>\n      .other-player {\n        border-radius: 50%;\n        width: 200px; \n        height: 200px; \n        background-color: #ffffffe0;\n        position: absolute; \n        transform-style: preserve-3d; \n        box-shadow: -20px -16px 8px 0 rgb(255 192 192), 20px 15px 20px 0 rgb(255 230 137);\n        filter: blur(20px);\n        animation-duration: 4s;\n        animation-name: pulse-player;\n        animation-iteration-count: infinite;\n      }\n\n      @keyframes pulse-player {\n        from { box-shadow: -20px -16px 8px 0 rgb(255 192 192), 20px 15px 20px 0 rgb(255 230 137);}\n        to { box-shadow: -15px -10px 8px 0 rgb(280 150 192), 20px 15px 20px 0 rgb(220 180 137);}\n      }\n    </style><div id=\"other_players\"></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func SceneTransform(player *gameobjects.Player) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templ.Raw("<style id='scene_transform'> #scene{"+
			fmt.Sprintf(`
  transform: perspective(4096px)
  rotate3d(0,1,0,%fdeg)
  rotate3d(1,0,0,%fdeg)
  rotate3d(0,0,1,%fdeg)
  translate3d(%fpx,%fpx,%fpx) scale3d(200,200,200); `,
				player.Rotation.Y,
				player.Rotation.X,
				player.Rotation.Z,
				player.Position.X*-200,
				player.Position.Y*200,
				player.Position.Z*-200,
			)+
			"}</style>").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
