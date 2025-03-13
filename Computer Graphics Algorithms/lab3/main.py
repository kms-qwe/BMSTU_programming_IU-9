import glfw
from OpenGL.GL import *
from OpenGL.GLU import *
import numpy as np
import math

def generate_twisted_pyramid_vertices(n_sides, n_slices, base_radius, height, twist_angle_total):

    vertices = []
    for i in range(n_slices):
        t = i / (n_slices - 1)
        scale = 1 - t
        twist = twist_angle_total * t
        slice_vertices = []
        for j in range(n_sides):
            angle = 2 * math.pi * j / n_sides + twist
            x = base_radius * scale * math.cos(angle)
            y = base_radius * scale * math.sin(angle)
            z = height * t
            slice_vertices.append((x, y, z))
        vertices.append(slice_vertices)
    return vertices

def draw_twisted_pyramid(n_sides=4, n_slices=20, base_radius=1.0, height=2.0, twist_angle_total=2*math.pi):

    vertices = generate_twisted_pyramid_vertices(n_sides, n_slices, base_radius, height, twist_angle_total)
    
    # Задаем цвет (например, белый) для отрисовки
    glColor3f(1, 1, 1)
    
    # Отрисовка боковых граней
    for i in range(n_slices - 1):
        if i == n_slices - 2:
            glBegin(GL_TRIANGLES)
            for j in range(n_sides):
                next_j = (j + 1) % n_sides
                glVertex3fv(vertices[i][j])
                glVertex3fv(vertices[i][next_j])
                glVertex3fv(vertices[i+1][j])
            glEnd()
        else:
            glBegin(GL_QUADS)
            for j in range(n_sides):
                next_j = (j + 1) % n_sides
                glVertex3fv(vertices[i][j])
                glVertex3fv(vertices[i+1][j])
                glVertex3fv(vertices[i+1][next_j])
                glVertex3fv(vertices[i][next_j])
            glEnd()
    
    # Отрисовка основания
    glBegin(GL_POLYGON)
    for j in range(n_sides):
        glVertex3fv(vertices[0][j])
    glEnd()

def main():
    if not glfw.init():
        return
    window = glfw.create_window(800, 600, "Скрученная пирамида", None, None)
    if not window:
        glfw.terminate()
        return
    glfw.make_context_current(window)
    
    glEnable(GL_DEPTH_TEST)
    
    # Настройка матрицы проекции
    glMatrixMode(GL_PROJECTION)
    glLoadIdentity()
    gluPerspective(45, 800/600, 0.1, 50.0)
    glMatrixMode(GL_MODELVIEW)
    
    # Задаем цвет фона (например, темно-серый, чтобы контрастировал с белой фигурой)
    glClearColor(0.1, 0.1, 0.1, 1)
    
    while not glfw.window_should_close(window):
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT)
        glLoadIdentity()
        # Настройка камеры
        gluLookAt(3, 3, 3, 0, 0, 0, 0, 0, 1)
        # Анимация вращения
        glRotatef(glfw.get_time() * 50, 0, 0, 1)
        
        # Отрисовка скрученной пирамиды
        draw_twisted_pyramid(n_sides=4, n_slices=20, base_radius=1.0, height=2.0, twist_angle_total=2*math.pi)
        
        glfw.swap_buffers(window)
        glfw.poll_events()
    glfw.terminate()

if __name__ == "__main__":
    main()
