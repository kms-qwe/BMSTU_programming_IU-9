import javax.swing.*;
import java.awt.*;
import java.awt.geom.AffineTransform;
import java.awt.geom.Path2D;

public class PictureForm extends JPanel {
    private double a;
    private double b;
    private double alphaDegrees;
    private double betaDegrees;

    public PictureForm(double a, double b, double alphaDegrees, double betaDegrees) {
        this.a = a;
        this.b = b;
        this.alphaDegrees = alphaDegrees;
        this.betaDegrees = betaDegrees;
    }

    @Override
    protected void paintComponent(Graphics g) {
        super.paintComponent(g);

        Graphics2D g2d = (Graphics2D) g;

        // Очищаем область рисования
        g2d.clearRect(0, 0, getWidth(), getHeight());

        // Создаем путь для треугольника
        Path2D path = new Path2D.Double();
        path.moveTo(0, 0);
        path.lineTo(a, 0);
        path.lineTo(0, b);
        path.closePath();

        // Поворачиваем треугольник в соответствии с заданными углами
        AffineTransform transform = AffineTransform.getRotateInstance(Math.toRadians(alphaDegrees), 0, 0);
        transform.concatenate(AffineTransform.getShearInstance(0, b / a));
        Shape transformedShape = transform.createTransformedShape(path);

        // Рисуем треугольник
        g2d.setColor(Color.blue);
        g2d.fill(transformedShape);
        g2d.setColor(Color.black);
        g2d.draw(transformedShape);
    }

    public static void main(String[] args) {
        SwingUtilities.invokeLater(() -> {
            double a = 100.0; // задаем катеты треугольника
            double b = 150.0;
            double alphaDegrees = 30.0; // задаем углы в градусах
            double betaDegrees = 60.0;

            JFrame frame = new JFrame("Triangle");
            frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
            frame.setSize(300, 300);
            frame.setLocationRelativeTo(null);

            PictureForm triangleForm = new PictureForm(a, b, alphaDegrees, betaDegrees);
            frame.add(triangleForm);

            frame.setVisible(true);
        });
    }
}
