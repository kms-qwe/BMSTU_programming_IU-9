import glfw
import numpy as np
import OpenGL.GL as gl
import math

rotation_angle = 0.0

def key_callback(window, key, scancode, action, mods):
    global rotation_angle
    if key == glfw.KEY_SPACE and action == glfw.PRESS:
        rotation_angle += 10  

def draw_dodecagon():
    global rotation_angle

    vertices = []
    num_sides = 12
    radius = 0.5

    for i in range(num_sides):
        angle = 2 * math.pi * i / num_sides
        x = radius * math.cos(angle)
        y = radius * math.sin(angle)
        vertices.append((x, y))

    vertices = np.array(vertices, dtype=np.float32)

    gl.glPushMatrix()
    gl.glRotatef(rotation_angle, 0, 0, 1)  

    gl.glBegin(gl.GL_POLYGON)
    gl.glColor3f(0.2, 0.6, 1.0)  
    for v in vertices:
        gl.glVertex2f(v[0], v[1])
    gl.glEnd()

    gl.glPopMatrix()

def main():
    if not glfw.init():
        return

    window = glfw.create_window(600, 600, "12-угольник", None, None)
    if not window:
        glfw.terminate()
        return

    glfw.make_context_current(window)
    glfw.set_key_callback(window, key_callback)

    while not glfw.window_should_close(window):
        gl.glClear(gl.GL_COLOR_BUFFER_BIT)
        gl.glLoadIdentity()

        draw_dodecagon()

        glfw.swap_buffers(window)
        glfw.poll_events()

    glfw.terminate()

if __name__ == "__main__":
    main()
