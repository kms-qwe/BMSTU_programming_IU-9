import glfw
from OpenGL.GL import *
from OpenGL.GLU import *
import numpy as np
import math
import time

# ----------------------------
# Глобальные параметры анимации
# ----------------------------
# Начальная позиция объекта (текущие координаты в пространстве)
pos = np.array([0.0, 0.0, 0.0], dtype=float)
# Начальная скорость объекта по осям X, Y, Z
vel = np.array([1.2, 0.9, 1.5], dtype=float)
# Половина размера ограничивающего куба (граница ±bounds по каждой оси)
bounds = 2.0

# ----------------------------
# Параметры источников света
# ----------------------------
# Массив словарей, каждый задаёт один источник:
# - id: идентификатор (GL_LIGHT0, GL_LIGHT1, ...)
# - ambient, diffuse, specular: компоненты цвета
# - pos: положение в мировых координатах
# - color: цвет сферы, показывающей источник
light_params = [
    {   # Красный источник света
        "id": GL_LIGHT0,
        "ambient": [0.2, 0.0, 0.0, 1.0],
        "diffuse": [1.0, 0.0, 0.0, 1.0],
        "specular":[1.0, 0.2, 0.2, 1.0],
        "pos":    [3.0, 3.0, 3.0, 1.0],
        "color":  (1.0, 0.0, 0.0)
    },
    {   # Синий источник света
        "id": GL_LIGHT1,
        "ambient": [0.0, 0.0, 0.2, 1.0],
        "diffuse": [0.0, 0.0, 1.0, 1.0],
        "specular":[0.2, 0.2, 1.0, 1.0],
        "pos":    [-3.0, -3.0, 2.0, 1.0],
        "color":  (0.0, 0.0, 1.0)
    }
]

def init_lighting_and_material():
    """
    Инициализация системы освещения и свойств материала.
    Вызывается один раз при старте приложения.
    """
    # Включаем расчёт освещения в OpenGL
    glEnable(GL_LIGHTING)

    # Устанавливаем глобальный фоновый (ambient) свет модели
    glLightModelfv(GL_LIGHT_MODEL_AMBIENT, [0.05, 0.05, 0.05, 1.0])

    # Настраиваем каждый источник из списка light_params
    for lp in light_params:
        glEnable(lp["id"])  # включаем конкретный источник
        # Устанавливаем компоненты освещения
        glLightfv(lp["id"], GL_AMBIENT,  lp["ambient"])
        glLightfv(lp["id"], GL_DIFFUSE,  lp["diffuse"])
        glLightfv(lp["id"], GL_SPECULAR, lp["specular"])
        # Позицию обновим каждый кадр после установки камеры
        glLightfv(lp["id"], GL_POSITION, lp["pos"])

    # Настройка свойств материала для всех отрисовываемых объектов
    glMaterialfv(GL_FRONT_AND_BACK, GL_AMBIENT,   [0.2, 0.2, 0.2, 1.0])
    glMaterialfv(GL_FRONT_AND_BACK, GL_DIFFUSE,   [0.8, 0.5, 0.3, 1.0])
    glMaterialfv(GL_FRONT_AND_BACK, GL_SPECULAR,  [1.0, 1.0, 1.0, 1.0])
    glMaterialf (GL_FRONT_AND_BACK, GL_SHININESS, 50.0)  # блеск (0–128)

def load_procedural_texture():
    """
    Генерация и загрузка процедурной 'шахматной' текстуры в OpenGL.
    Используется для управления интенсивностью поверхности (GL_LUMINANCE).
    """
    glEnable(GL_TEXTURE_2D)  # включаем 2D-текстурирование
    glTexEnvf(GL_TEXTURE_ENV, GL_TEXTURE_ENV_MODE, GL_MODULATE)  # умножаем текстуру на цвет материала

    tex_size = 64  # размер текстуры в пикселях
    # Создаем массив оттенков серого (0–255)
    checker = np.zeros((tex_size, tex_size), dtype=np.uint8)
    for i in range(tex_size):
        for j in range(tex_size):
            # рисуем клетки размера 8×8 по принципу 'шахматка'
            checker[i,j] = 255 if (i//8 + j//8) % 2 == 0 else 128

    # Загружаем данные в виде LUMINANCE-канала
    glTexImage2D(GL_TEXTURE_2D, 0, GL_LUMINANCE, tex_size, tex_size, 0,
                 GL_LUMINANCE, GL_UNSIGNED_BYTE, checker)
    # Фильтрация: без сглаживания (пикселизация)
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST)
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST)

def generate_twisted_pyramid_vertices(n_sides, n_slices, base_radius, height, twist_angle_total):
    """
    Генерация вершин скрученной пирамиды.
    - n_sides: число граней основания
    - n_slices: число 'слоёв' по высоте
    - base_radius: радиус основания
    - height: общая высота
    - twist_angle_total: общий угол закрутки (рад)
    Возвращает список списков вершин: verts[slice][side] = (x,y,z)
    """
    verts = []
    for i in range(n_slices):
        t = i / (n_slices - 1)    # доля от основания к вершине (0..1)
        scale = 1 - t             # масштаб радиуса на этом уровне
        twist = twist_angle_total * t  # текущая закрутка
        slice_verts = []
        for j in range(n_sides):
            angle = 2*math.pi*j/n_sides + twist
            x = base_radius*scale*math.cos(angle)
            y = base_radius*scale*math.sin(angle)
            z = height * t
            slice_verts.append((x, y, z))
        verts.append(slice_verts)
    return verts

def draw_twisted_pyramid(n_sides=4, n_slices=20, base_radius=1.0, height=2.0, twist_angle_total=2*math.pi):
    """
    Отрисовка скрученной пирамиды по ранее сгенерированным вершинам.
    Используем QUADS для боковых поверхностей и TRIANGLES для вершины.
    """
    verts = generate_twisted_pyramid_vertices(n_sides, n_slices, base_radius, height, twist_angle_total)

    for i in range(n_slices - 1):
        # Последний срез — треугольники, остальные — квадраты
        primitive = GL_TRIANGLES if (i == n_slices - 2) else GL_QUADS
        glBegin(primitive)
        for j in range(n_sides):
            next_j = (j + 1) % n_sides

            # Вычисление нормали для данной грани (flat shading)
            v0 = np.array(verts[i][j])
            v1 = np.array(verts[i+1][j])
            v2 = np.array(verts[i][next_j])
            normal = np.cross(v1 - v0, v2 - v0)
            normal = normal / np.linalg.norm(normal)
            glNormal3fv(normal)

            # Текстурные координаты (поверхность растягивается по всей пирамиде)
            glTexCoord2f(j/(n_sides-1), i/(n_slices-1))
            glVertex3fv(verts[i][j])

            if primitive == GL_QUADS:
                # 2-я, 3-я, 4-я точки квадрата
                glTexCoord2f(j/(n_sides-1), (i+1)/(n_slices-1)); glVertex3fv(verts[i+1][j])
                glTexCoord2f((j+1)/(n_sides-1), (i+1)/(n_slices-1)); glVertex3fv(verts[i+1][next_j])
                glTexCoord2f((j+1)/(n_sides-1), i/(n_slices-1)); glVertex3fv(verts[i][next_j])
            else:
                # для треугольника: две дополнительные точки
                glTexCoord2f((j+1)/(n_sides-1), (i+1)/(n_slices-1)); glVertex3fv(verts[i][next_j])
                glTexCoord2f(j/(n_sides-1), (i+1)/(n_slices-1)); glVertex3fv(verts[i+1][j])
        glEnd()

    # Отрисовка основания пирамиды (полигон)
    glBegin(GL_POLYGON)
    glNormal3f(0, 0, -1)  # нормаль вниз
    for j in range(n_sides):
        glTexCoord2f(j/(n_sides-1), 0)
        glVertex3fv(verts[0][j])
    glEnd()

def draw_bounding_box(size=2.0, alpha=0.3):
    """
    Отрисовка полупрозрачного каркаса ограничивающего куба.
    - size: половина длины ребра
    - alpha: прозрачность (0.0–1.0)
    """
    # Чтобы линии были видны, временно отключаем освещение
    glDisable(GL_LIGHTING)
    glColor4f(1, 1, 1, alpha)  # белый цвет с заданной прозрачностью
    glLineWidth(2.0)           # толщина линий

    glBegin(GL_LINES)
    # 12 рёбер куба: перебираем пары точек
    for x in (-size, size):
        for y in (-size, size):
            glVertex3f(x, y, -size); glVertex3f(x, y,  size)
    for x in (-size, size):
        for z in (-size, size):
            glVertex3f(x, -size, z); glVertex3f(x,  size, z)
    for y in (-size, size):
        for z in (-size, size):
            glVertex3f(-size, y, z); glVertex3f( size, y, z)
    glEnd()

    # Восстанавливаем расчёт освещения для остальных объектов
    glEnable(GL_LIGHTING)

def draw_light_sources():
    """
    Отрисовка визуализации самих источников света в виде маленьких сфер.
    Сферы рисуются без освещения (чтобы быть всегда яркими).
    """
    glDisable(GL_LIGHTING)
    quad = gluNewQuadric()
    for lp in light_params:
        glPushMatrix()
        # Перемещаемся в позицию источника
        glTranslatef(*lp["pos"][:3])
        # Устанавливаем цвет сферы по diffuse-компоненте
        glColor3f(*lp["color"])
        # Рисуем сферу радиусом 0.1, 16×16 сегментов
        gluSphere(quad, 0.1, 16, 16)
        glPopMatrix()
    glEnable(GL_LIGHTING)

def main():
    # Инициализация GLFW
    if not glfw.init():
        return
    # Создание окна
    win = glfw.create_window(800, 600, "Лабораторная 6 – цветные источники", None, None)
    if not win:
        glfw.terminate()
        return
    glfw.make_context_current(win)

    # Включаем тест глубины, чтобы объекты правильно перекрывались
    glEnable(GL_DEPTH_TEST)

    # Настройка проекции (матрица проекции)
    glMatrixMode(GL_PROJECTION)
    glLoadIdentity()
    gluPerspective(45, 800/600, 0.1, 50.0)  # FOV=45°, aspect, near=0.1, far=50
    glMatrixMode(GL_MODELVIEW)

    # Фон сцены — тёмно-серый
    glClearColor(0.1, 0.1, 0.1, 1.0)

    # Инициализируем свет и материалы
    init_lighting_and_material()
    # Генерируем и загружаем текстуру
    load_procedural_texture()

    last_t = time.time()  # время предыдущего кадра
    # Главный цикл отрисовки
    while not glfw.window_should_close(win):
        t = time.time()
        dt = t - last_t
        last_t = t

        # Обновление позиции объекта: pos = pos + vel * dt
        global pos, vel
        pos += vel * dt
        # Проверка и отражение от стенок ограничивающего куба
        for i in range(3):
            if abs(pos[i]) > bounds:
                pos[i] = np.sign(pos[i]) * bounds  # ставим на границу
                vel[i] = -vel[i]                   # меняем направление скорости

        # Очистка буферов цвета и глубины
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT)
        glLoadIdentity()
        # Устанавливаем камеру: смотри из (5,5,5) в центр, 'вверх' по Z
        gluLookAt(5, 5, 5,  0, 0, 0,  0, 0, 1)

        # После установки вида обновляем позиции источников
        for lp in light_params:
            glLightfv(lp["id"], GL_POSITION, lp["pos"])

        # Рисуем полупрозрачный каркас куба
        draw_bounding_box(size=bounds, alpha=0.3)
        # Визуализируем сами источники света
        draw_light_sources()

        # Применяем трансформацию к движущейся пирамиде
        glTranslatef(*pos)
        glRotatef((t * 30) % 360, 0, 0, 1)  # вращаем вокруг Z со скоростью 30°/с

        # Наконец, рисуем скрученную пирамиду
        draw_twisted_pyramid()

        # Показываем результат на экране
        glfw.swap_buffers(win)
        glfw.poll_events()

    # Завершаем работу GLFW
    glfw.terminate()

if __name__ == "__main__":
    main()
