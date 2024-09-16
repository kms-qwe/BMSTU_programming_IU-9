import javax.swing.*;
import java.awt.*;
import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;

public class StarDrawingApp extends JFrame {

    private StarPanel starPanel;
    private JSlider verticesSlider;
    private JCheckBox fillCheckBox;

    public StarDrawingApp() {
        setTitle("N-Pointed Star Drawer");
        setSize(600, 600);
        setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        setLocationRelativeTo(null);

        starPanel = new StarPanel();
        verticesSlider = new JSlider(JSlider.HORIZONTAL, 3, 20, 3);
        fillCheckBox = new JCheckBox("Fill");

        verticesSlider.addChangeListener(e -> starPanel.setVertices(verticesSlider.getValue()));
        fillCheckBox.addActionListener(e -> starPanel.setFill(fillCheckBox.isSelected()));

        JPanel controlPanel = new JPanel();
        controlPanel.add(new JLabel("Vertices:"));
        controlPanel.add(verticesSlider);
        controlPanel.add(fillCheckBox);

        add(starPanel, BorderLayout.CENTER);
        add(controlPanel, BorderLayout.SOUTH);
    }

    public static void main(String[] args) {
        SwingUtilities.invokeLater(() -> {
            StarDrawingApp app = new StarDrawingApp();
            app.setVisible(true);
        });
    }
}

class StarPanel extends JPanel {

    private int vertices = 5;
    private boolean fill = false;

    public void setVertices(int vertices) {
        this.vertices = vertices;
        repaint();
    }

    public void setFill(boolean fill) {
        this.fill = fill;
        repaint();
    }

    @Override
    protected void paintComponent(Graphics g) {
        super.paintComponent(g);
        drawStar(g, getWidth() / 2, getHeight() / 2, Math.min(getWidth(), getHeight()) / 3, vertices, fill);
    }

    private void drawStar(Graphics g, int xCenter, int yCenter, int radius, int vertices, boolean fill) {
        double angle = Math.PI / vertices;
        int[] xPoints = new int[2 * vertices];
        int[] yPoints = new int[2 * vertices];

        for (int i = 0; i < 2 * vertices; i++) {
            int r = (i % 2 == 0) ? radius : radius / 2;
            xPoints[i] = xCenter + (int) (r * Math.sin(i * angle));
            yPoints[i] = yCenter - (int) (r * Math.cos(i * angle));
        }

        if (fill) {
            g.fillPolygon(xPoints, yPoints, 2 * vertices);
        } else {
            g.drawPolygon(xPoints, yPoints, 2 * vertices);
        }
    }
}
