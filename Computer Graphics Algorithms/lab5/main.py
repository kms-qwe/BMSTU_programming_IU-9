import glfw
from OpenGL.GL import *
import sys

# Region codes
INSIDE = 0  # 0000
LEFT   = 1  # 0001
RIGHT  = 2  # 0010
BOTTOM = 4  # 0100
TOP    = 8  # 1000

# Clipping rectangle (initialized later)
xmin, ymin, xmax, ymax = 0.0, 0.0, 0.0, 0.0

# Store user-defined primitives
lines = []       # list of (x0,y0,x1,y1)
clipped_lines = []
rect_defined = False
points = []      # temporary points for input

# Compute region code for a point (x,y)
def compute_code(x, y):
    code = INSIDE
    if x < xmin:
        code |= LEFT
    elif x > xmax:
        code |= RIGHT
    if y < ymin:
        code |= BOTTOM
    elif y > ymax:
        code |= TOP
    return code

# Cohen–Sutherland clipping algorithm

def cohen_sutherland_clip(x0, y0, x1, y1):
    code0 = compute_code(x0, y0)
    code1 = compute_code(x1, y1)
    accept = False

    while True:
        if code0 == 0 and code1 == 0:
            accept = True
            break
        elif (code0 & code1) != 0:
            break
        else:
            out_code = code0 or code1
            if out_code & TOP:
                x = x0 + (x1 - x0) * (ymax - y0) / (y1 - y0)
                y = ymax
            elif out_code & BOTTOM:
                x = x0 + (x1 - x0) * (ymin - y0) / (y1 - y0)
                y = ymin
            elif out_code & RIGHT:
                y = y0 + (y1 - y0) * (xmax - x0) / (x1 - x0)
                x = xmax
            else:  # LEFT
                y = y0 + (y1 - y0) * (xmin - x0) / (x1 - x0)
                x = xmin

            if out_code == code0:
                x0, y0 = x, y
                code0 = compute_code(x0, y0)
            else:
                x1, y1 = x, y
                code1 = compute_code(x1, y1)

    return (x0, y0, x1, y1) if accept else None

# Process click in normalized coords

def process_click(window, x, y):
    global points, rect_defined, xmin, ymin, xmax, ymax
    # Normalize using window size (points)
    w, h = glfw.get_window_size(window)
    nx =  (x / w) * 2 - 1
    ny = 1 - (y / h) * 2
    points.append((nx, ny))
    if not rect_defined and len(points) == 2:
        (x0, y0), (x1, y1) = points
        xmin, xmax = min(x0, x1), max(x0, x1)
        ymin, ymax = min(y0, y1), max(y0, y1)
        rect_defined = True
        points.clear()
    elif rect_defined and len(points) == 2:
        (x0, y0), (x1, y1) = points
        lines.append((x0, y0, x1, y1))
        clipped = cohen_sutherland_clip(x0, y0, x1, y1)
        if clipped:
            clipped_lines.append(clipped)
        points.clear()

# GLFW callbacks

def mouse_button_callback(window, button, action, mods):
    if button == glfw.MOUSE_BUTTON_LEFT and action == glfw.PRESS:
        x, y = glfw.get_cursor_pos(window)
        process_click(window, x, y)


def framebuffer_size_callback(window, width, height):
    glViewport(0, 0, width, height)


def key_callback(window, key, scancode, action, mods):
    if key == glfw.KEY_ESCAPE and action == glfw.PRESS:
        glfw.set_window_should_close(window, True)


def draw_scene():
    if rect_defined:
        glColor3f(0,1,0)
        glBegin(GL_LINE_LOOP)
        glVertex2f(xmin, ymin)
        glVertex2f(xmax, ymin)
        glVertex2f(xmax, ymax)
        glVertex2f(xmin, ymax)
        glEnd()
    glColor3f(1,1,1)
    glBegin(GL_LINES)
    for x0, y0, x1, y1 in lines:
        glVertex2f(x0, y0)
        glVertex2f(x1, y1)
    glEnd()
    glColor3f(1,0,0)
    glBegin(GL_LINES)
    for cx0, cy0, cx1, cy1 in clipped_lines:
        glVertex2f(cx0, cy0)
        glVertex2f(cx1, cy1)
    glEnd()


def main():
    if not glfw.init():
        print("Failed to initialize GLFW")
        return
    # Retina support on macOS
    glfw.window_hint(glfw.COCOA_RETINA_FRAMEBUFFER, glfw.TRUE)
    width, height = 800, 600
    window = glfw.create_window(width, height, "Cohen-Sutherland Clipping", None, None)
    if not window:
        glfw.terminate()
        print("Failed to create GLFW window")
        return

    glfw.make_context_current(window)
    # Setup viewport for high-DPI
    fb_w, fb_h = glfw.get_framebuffer_size(window)
    glViewport(0, 0, fb_w, fb_h)
    glfw.set_framebuffer_size_callback(window, framebuffer_size_callback)
    glfw.set_mouse_button_callback(window, mouse_button_callback)
    glfw.set_key_callback(window, key_callback)

    while not glfw.window_should_close(window):
        glfw.poll_events()
        glClear(GL_COLOR_BUFFER_BIT)
        draw_scene()
        glfw.swap_buffers(window)

    glfw.terminate()

if __name__ == '__main__':
    main()
