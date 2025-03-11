import glfw
from OpenGL.GL import *
import numpy as np

# Параметры вращения
angle_x, angle_y = 0, 0
wireframe = False  # Переключение режима

def apply_isometric_projection():
    # Матрица изометрической проекции
    iso_matrix = np.array([
        [0.707, 0, -0.707, 0],  # Поворот на 45 градусов вокруг Y
        [0.408, 0.816, 0.408, 0],  # Поворот на 35.264 градусов вокруг X
        [0.577, -0.577, 0.577, 0],
        [0, 0, 0, 1]
    ], dtype=np.float32)
    
    glMultMatrixf(iso_matrix.T.flatten())

def main():
    if not glfw.init():
        return
    window = glfw.create_window(800, 800, "Isometric Cube", None, None)
    if not window:
        glfw.terminate()
        return
    glfw.make_context_current(window)
    glfw.set_key_callback(window, key_callback)
    
    glEnable(GL_DEPTH_TEST)
    glClearColor(0.2, 0.2, 0.2, 1.0)  # Измененный цвет фона

    while not glfw.window_should_close(window):
        display()
        glfw.swap_buffers(window)
        glfw.poll_events()
    
    glfw.destroy_window(window)
    glfw.terminate()

def draw_cube():
    # Новые цвета граней для различимости
    face_colors = [
        [0.8, 0.0, 0.0],  # Красный
        [0.0, 0.8, 0.0],  # Зеленый
        [0.0, 0.0, 0.8],  # Синий
        [0.8, 0.8, 0.0],  # Желтый
        [0.8, 0.0, 0.8],  # Фиолетовый
        [0.0, 0.8, 0.8],  # Голубой
    ]

    vertices = [
        [-0.25, -0.25, -0.25], [0.25, -0.25, -0.25], [0.25, 0.25, -0.25], [-0.25, 0.25, -0.25],
        [-0.25, -0.25, 0.25], [0.25, -0.25, 0.25], [0.25, 0.25, 0.25], [-0.25, 0.25, 0.25]
    ]
    faces = [
        [0, 1, 2, 3], [4, 5, 6, 7], [0, 4, 7, 3],
        [1, 5, 6, 2], [3, 2, 6, 7], [0, 1, 5, 4]
    ]
    
    glPolygonMode(GL_FRONT_AND_BACK, GL_LINE if wireframe else GL_FILL)
    
    for i, face in enumerate(faces):
        glBegin(GL_QUADS)
        glColor3fv(face_colors[i])
        for vertex in face:
            glVertex3fv(vertices[vertex])
        glEnd()

def display():
    global angle_x, angle_y
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT)
    
    # Вращающийся куб (передний)
    glLoadIdentity()
    apply_isometric_projection()
    glTranslatef(-0.8, 0, 0)  # Смещаем его левее
    glRotatef(angle_x, 1, 0, 0)
    glRotatef(angle_y, 0, 1, 0)
    draw_cube()
    
    # Статичный куб (задний)
    glLoadIdentity()
    apply_isometric_projection()
    glTranslatef(0.8, 0, 0)  # Смещаем его правее
    draw_cube()

def key_callback(window, key, scancode, action, mods):
    global angle_x, angle_y, wireframe
    if action in [glfw.PRESS, glfw.REPEAT]:
        if key == glfw.KEY_RIGHT:
            angle_y -= 3
        if key == glfw.KEY_LEFT:
            angle_y += 3
        if key == glfw.KEY_UP:
            angle_x -= 3
        if key == glfw.KEY_DOWN:
            angle_x += 3
        if key == glfw.KEY_M:
            wireframe = not wireframe
        if key == glfw.KEY_ESCAPE:
            glfw.set_window_should_close(window, True)

if __name__ == "__main__":
    main()
