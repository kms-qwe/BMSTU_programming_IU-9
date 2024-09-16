import javax.swing.*;
import java.awt.*;
import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;

public class TriangleGUIForm extends JFrame {
    private JTextField angleTextField;
    private TrianglePanel trianglePanel;

    public TriangleGUIForm() {
        setTitle("Прямоугольный треугольник");
        setSize(400, 400);
        setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        setLocationRelativeTo(null);

        JPanel controlPanel = new JPanel();
        JLabel angleLabel = new JLabel("Введите значение острого угла (в градусах):");
        angleTextField = new JTextField(10);
        angleTextField.addActionListener(new ActionListener() {
            @Override
            public void actionPerformed(ActionEvent e) {
                drawTriangle();
            }
        });
        JButton drawButton = new JButton("Нарисовать");
        drawButton.addActionListener(new ActionListener() {
            @Override
            public void actionPerformed(ActionEvent e) {
                drawTriangle();
            }
        });
        controlPanel.add(angleLabel);
        controlPanel.add(angleTextField);
        controlPanel.add(drawButton);
        add(controlPanel, BorderLayout.NORTH);

        trianglePanel = new TrianglePanel();
        add(trianglePanel, BorderLayout.CENTER);
    }

    private void drawTriangle() {
            double angle = Double.parseDouble(angleTextField.getText());
            trianglePanel.setAngle(angle);
            trianglePanel.repaint();
    }

    public static void main(String[] args) {
        SwingUtilities.invokeLater(new Runnable() {
            @Override
            public void run() {
                TriangleGUIForm triangleGUIForm = new TriangleGUIForm();
                triangleGUIForm.setVisible(true);
            }
        });
    }
}

class TrianglePanel extends JPanel {
    private double angle;

    public TrianglePanel() {
        angle = 45.0;
    }

    public void setAngle(double angle) {
        this.angle = angle;
    }

    @Override
    protected void paintComponent(Graphics g) {
        super.paintComponent(g);

        int width = getWidth();
        int height = getHeight();
        int x1 = 300;
        int y1 = height - 300;
        int x2 = x1 + 200;
        int y2 = y1;
        int x3 = x1;
        int y3 = y1 - (int) (Math.tan(Math.toRadians(angle)) * 200);

        Graphics2D g2d = (Graphics2D) g;
        g2d.setColor(Color.BLUE);
        g2d.drawLine(x1, y1, x2, y2);
        g2d.drawLine(x1, y1, x3, y3);
        g2d.drawLine(x2, y2, x3, y3);
    }
}
