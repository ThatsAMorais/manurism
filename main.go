package main

import (
    "fmt"
    _ "image/png"
//    "log"
    "runtime"
    "github.com/ThatsAMorais/gogl-engine/game"

    "github.com/go-gl/glfw/v3.1/glfw"
//    "github.com/go-gl/mathgl/mgl32"
    "github.com/go-gl/mathgl/mgl32"
//    "github.com/d4l3k/go-pry/pry"
    "github.com/go-gl/gl/v3.3-core/gl"
    "log"
)

const WindowWidth = 1024
const WindowHeight = 768

type Movement int
const (
    Forward     int = 1 << 0
    Reverse     int = 1 << 1
    StrafeLeft  int = 1 << 2
    StrafeRight int = 1 << 3
)

var prevXPos = 0.0
var prevYPos = 0.0
var vX, vY, vZ = 0.0, 0.0, 0.0
var inputs int

func init() {
    // GLFW event handling must run on the main OS thread
    runtime.LockOSThread()
//    log.Println("Initializing glfw")
//    if err := glfw.Init(); err != nil {
//        log.Fatalln("failed to initialize glfw:", err)
//    }
}

func main() {
    /**
     * Main method for the game
     *
     */

    gameObj1 := game.NewGameObject()
    gameObj1.Pos = mgl32.Vec3{float32(0.5), float32(0.0), float32(-0.5)}
    gameObj1.Velocity = mgl32.Vec3{float32(0.5), float32(0.0), float32(-0.5)}
    gameObj1.Name = "Player"
    gameObj2 := game.NewGameObject()
    gameObj2.Pos = mgl32.Vec3{float32(0.5), float32(0.0), float32(-0.5)}
    gameObj1.Velocity = mgl32.Vec3{float32(0.5), float32(0.0), float32(-0.5)}
    gameObj2.Name = "Obj 1"
    gameObjects := []game.GameObject{*gameObj1, *gameObj2}

    scene := game.NewScene()
    scene.GameObjects = gameObjects
    scene.Renderer = game.Renderer{}

    _game := game.NewGame(*scene, keyPressHandler, mouseMovementHandler)

    if err := glfw.Init(); err != nil {
        log.Fatalln("failed to initialize glfw:", err)
    }

    glfw.WindowHint(glfw.Resizable, glfw.False)
    glfw.WindowHint(glfw.ContextVersionMajor, 3)
    glfw.WindowHint(glfw.ContextVersionMinor, 3)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
    window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Game", nil, nil)
    if err != nil {
        panic(err)
    }
    window.MakeContextCurrent()

    // Initialize Glow
    if err := gl.Init(); err != nil {
        panic(err)
    }

    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Println("OpenGL version", version)

    gl.Viewport(0, 0, WindowWidth, WindowHeight)

    // Configure the vertex and fragment shaders.
    shader := game.NewShader()
    program, err := shader.GetProgram()  // newProgram(vertexShaderSource, fragmentShaderSource)
    if err != nil {
        panic(err)
    }
    gl.UseProgram(program)

    // Register pre-game-loop handlers
    window.SetKeyCallback(keyPressHandler)
    window.SetCursorPosCallback(mouseMovementHandler)
    window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
    defer window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)

    _game.Initialize()
    _game.SetWindow(window)
    defer _game.Stop()
    for !window.ShouldClose() {
        var verticies = []float32{
            -.5, -.5, 0.0, 0.0, 1.0,
            .5, -.5, 0.0, 1.0, 0.0,
            -.5, .5, 0.0, 0.0, 1.0,
        }

        var vao uint32
        gl.GenVertexArrays(1, &vao)
        gl.BindVertexArray(vao)

        var vbo uint32
        gl.GenBuffers(1, &vbo)
        gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
        gl.BufferData(gl.ARRAY_BUFFER, len(verticies)*4, gl.Ptr(verticies), gl.STATIC_DRAW)

        vertAttrib := uint32(gl.GetAttribLocation(_game.GetProgram(), gl.Str("vert\x00")))
        gl.EnableVertexAttribArray(vertAttrib)
        gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

        texCoordAttrib := uint32(gl.GetAttribLocation(_game.GetProgram(), gl.Str("vertTexCoord\x00")))
        gl.EnableVertexAttribArray(texCoordAttrib)
        gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))


        frameCount := 0
        previousTime := glfw.GetTime()
        // Update time tracks
        if frameCount % 240 == 0 {
            time := glfw.GetTime()
            elapsed := time - previousTime
            previousTime = time
            fmt.Println("Frame count", frameCount, "Time for last 240 frames ", elapsed)
        }

        timeStep, relativeVel := _game.ParseInputs()


        for i := 0; i < len(_game.MyScene.GameObjects); i++ {
            _game.MyScene.GameObjects[i].UpdatePosition(timeStep, relativeVel)
        }

        _game.MyScene.Render()

        // Handle events
        glfw.PollEvents()

        // Clear screen
        gl.ClearColor(float32(0.2), float32(0.2), float32(0.2), float32(1.0))
        gl.Clear(gl.COLOR_BUFFER_BIT)

        gl.BindVertexArray(vao)
        gl.DrawArrays(gl.TRIANGLES, 0, 3)

        //g.window.SwapBuffers()
        frameCount++
        window.SwapBuffers()
    }
}

func mouseMovementHandler(w *glfw.Window, xpos float64, ypos float64) {
    dx, dy := prevXPos - xpos, prevYPos - ypos
    fmt.Printf("New cursor position is (%v, %v). Cursor moved (%v, %v)\n", xpos, ypos, dx, dy)
    prevXPos, prevYPos = xpos, ypos
}

func keyPressHandler(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
    if action == glfw.Press {
        if key == glfw.KeyEscape {
            w.SetShouldClose(true)
        } else if key == glfw.KeyW {
            inputs = inputs | Forward
            fmt.Println("Starting to move forward now", inputs)
        } else if key == glfw.KeyA {
            inputs = inputs | StrafeLeft
            fmt.Println("Starting to move left now", inputs)
        } else if key == glfw.KeyS {
            inputs = inputs | StrafeRight
            fmt.Println("Starting to move back now", inputs)
        } else if key == glfw.KeyD {
            inputs = inputs | Reverse
            fmt.Println("Starting to move right now", inputs)
        }
    } else if action == glfw.Release {
        if key == glfw.KeyW {
            inputs = inputs &^ Forward
            fmt.Println("Done moving forward now", inputs)
        } else if key == glfw.KeyA {
            inputs = inputs &^ StrafeLeft
            fmt.Println("Done moving left now", inputs)
        } else if key == glfw.KeyS {
            inputs = inputs &^ StrafeRight
            fmt.Println("Done moving back now", inputs)
        } else if key == glfw.KeyD {
            inputs = inputs &^ Reverse
            fmt.Println("Done moving right now", inputs)
        }
    }
}
