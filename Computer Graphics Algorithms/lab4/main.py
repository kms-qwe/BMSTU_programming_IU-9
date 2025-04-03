import glfw
from OpenGL.GL import *
import numpy as np
import math

WINDOW_WIDTH = 800
WINDOW_HEIGHT = 600

pixel_buffer = np.zeros((WINDOW_HEIGHT, WINDOW_WIDTH, 3), dtype=np.uint8)

points = []

def set_pixel(buffer, x, y, color):
    if 0 <= x < buffer.shape[1] and 0 <= y < buffer.shape[0]:
        buffer[y, x] = color

def bresenham_line(buffer, x0, y0, x1, y1, color):
    dx = abs(x1 - x0)
    dy = abs(y1 - y0)
    sx = 1 if x0 < x1 else -1
    sy = 1 if y0 < y1 else -1
    err = dx - dy

    while True:
        set_pixel(buffer, x0, y0, color)
        if x0 == x1 and y0 == y1:
            break
        e2 = 2 * err
        if e2 > -dy:
            err -= dy
            x0 += sx
        if e2 < dx:
            err += dx
            y0 += sy

def draw_polygon_edges(buffer, points, color):
    n = len(points)
    for i in range(n):
        x0, y0 = points[i]
        x1, y1 = points[(i + 1) % n]
        bresenham_line(buffer, int(x0), int(y0), int(x1), int(y1), color)

def fill_polygon_scanline(buffer, points, fill_color):
    if not points:
        return

    ys = [p[1] for p in points]
    y_min = int(min(ys))
    y_max = int(max(ys))
    n = len(points)

    for y in range(y_min, y_max + 1):
        intersections = []
        # По всем ребрам
        for i in range(n):
            p1 = points[i]
            p2 = points[(i + 1) % n]
            if p1[1] == p2[1]:
                continue
            if (y >= min(p1[1], p2[1])) and (y < max(p1[1], p2[1])):
                x_int = p1[0] + (y - p1[1]) * (p2[0] - p1[0]) / (p2[1] - p1[1])
                intersections.append(x_int)
        intersections.sort()
        for i in range(0, len(intersections), 2):
            if i + 1 < len(intersections):
                x_start = int(math.ceil(intersections[i]))
                x_end = int(math.floor(intersections[i + 1]))
                for x in range(x_start, x_end + 1):
                    set_pixel(buffer, x, y, fill_color)

def apply_filter(buffer):
    kernel = np.array([[1, 2, 1],
                       [2, 4, 2],
                       [1, 2, 1]], dtype=np.float32)
    kernel = kernel / kernel.sum()
    filtered = np.copy(buffer)
    h, w, _ = buffer.shape
    for y in range(1, h - 1):
        for x in range(1, w - 1):
            region = buffer[y - 1:y + 2, x - 1:x + 2]
            for c in range(3):
                filtered[y, x, c] = np.clip(np.sum(region[:, :, c] * kernel), 0, 255)
    return filtered

def resize_buffer(old_buffer, new_height, new_width):
    old_height, old_width, _ = old_buffer.shape
    if old_height == 0 or old_width == 0:
        return np.zeros((new_height, new_width, 3), dtype=np.uint8)
    new_y_indices = (np.linspace(0, old_height - 1, new_height)).astype(np.int32)
    new_x_indices = (np.linspace(0, old_width - 1, new_width)).astype(np.int32)
    new_buffer = old_buffer[new_y_indices][:, new_x_indices]
    return new_buffer

def key_callback(window, key, scancode, action, mods):
    global points, pixel_buffer
    if action == glfw.PRESS:
        if key == glfw.KEY_C:
            points = []
            pixel_buffer.fill(0)
        elif key == glfw.KEY_F:
            pixel_buffer[:] = apply_filter(pixel_buffer)

def mouse_button_callback(window, button, action, mods):
    global points, pixel_buffer
    if button == glfw.MOUSE_BUTTON_LEFT and action == glfw.PRESS:
        x, y = glfw.get_cursor_pos(window)
        y = WINDOW_HEIGHT - int(y)
        points.append((int(x), int(y)))
        if len(points) > 1:
            bresenham_line(pixel_buffer, points[-2][0], points[-2][1],
                           points[-1][0], points[-1][1], (255, 255, 255))
    elif button == glfw.MOUSE_BUTTON_RIGHT and action == glfw.PRESS:
        if len(points) > 2:
            bresenham_line(pixel_buffer, points[-1][0], points[-1][1],
                           points[0][0], points[0][1], (255, 255, 255))
            fill_polygon_scanline(pixel_buffer, points, (0, 255, 0))
            points.clear()

def framebuffer_size_callback(window, width, height):
    global WINDOW_WIDTH, WINDOW_HEIGHT, pixel_buffer
    old_buffer = pixel_buffer.copy()
    WINDOW_WIDTH = width
    WINDOW_HEIGHT = height
    glViewport(0, 0, width, height)
    pixel_buffer = resize_buffer(old_buffer, height, width)

def main():
    if not glfw.init():
        return

    window = glfw.create_window(WINDOW_WIDTH, WINDOW_HEIGHT, "Rasterization Demo", None, None)
    if not window:
        glfw.terminate()
        return

    glfw.make_context_current(window)
    glfw.set_key_callback(window, key_callback)
    glfw.set_mouse_button_callback(window, mouse_button_callback)
    glfw.set_framebuffer_size_callback(window, framebuffer_size_callback)

    glEnable(GL_TEXTURE_2D)
    texture = glGenTextures(1)
    
    while not glfw.window_should_close(window):
        glfw.poll_events()

        glBindTexture(GL_TEXTURE_2D, texture)
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST)
        glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST)
        glTexImage2D(GL_TEXTURE_2D, 0, GL_RGB, pixel_buffer.shape[1], pixel_buffer.shape[0],
                     0, GL_RGB, GL_UNSIGNED_BYTE, pixel_buffer)

        glClear(GL_COLOR_BUFFER_BIT)
        glLoadIdentity()

        glBegin(GL_QUADS)
        glTexCoord2f(0, 0); glVertex2f(-1, -1)
        glTexCoord2f(1, 0); glVertex2f(1, -1)
        glTexCoord2f(1, 1); glVertex2f(1, 1)
        glTexCoord2f(0, 1); glVertex2f(-1, 1)
        glEnd()

        glfw.swap_buffers(window)

    glfw.terminate()

if __name__ == "__main__":
    main()
