// THIS VERSION HAS A SUCCESSFUL CYCLE. NO COLOR CHANGE THOUGH 

function createPieChart(canvasId, labels, data, options = {}) {
  const canvas = document.getElementById(canvasId);
  if (!canvas) {
    console.error(`Canvas with id "${canvasId}" not found`);
    return null;
  }
  const ctx = canvas.getContext('2d');

  // destroy old chart on this canvas if it exists
  const existing = Chart.getChart(ctx);
  if (existing) existing.destroy();

  return new Chart(ctx, {
    type: 'doughnut',
    data: {
      labels: labels,
      datasets: [{
        data: data,
        backgroundColor: labels.map((_, i) => options.colors?.[i] || "#2B2D30"),
        borderWidth: 0
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      cutout: options.cutout || '80%',
      plugins: {
        legend: { display: false },
        ...options.plugins
      }
    },
    plugins: options.pluginsArray || []
  });
}

// === FETCH DATA ===
fetch("/api/monthly-totals")
  .then(res => res.json())
  .then(data => {
    const month = "Jan"; // can be dynamic
    const categories = data[month].Totals.map(e => e.Category);
    const totals = data[month].Totals.map(e => e.Total);
    const totalSum = totals.reduce((a, b) => a + b, 0);

    // === STATIC PIE ===
    const centerTextStatic = {
      id: 'centerTextStatic',
      afterDraw(chart) {
        const { ctx, chartArea: { width } } = chart;
        ctx.save();
        ctx.font = "bold 18px Saira";
        ctx.fillStyle = "#FFFFFF";
        ctx.textAlign = "center";
        ctx.textBaseline = "middle";
        ctx.fillText(`Spending by Category (${month})`, width / 2, 20);
        ctx.restore();
      }
    };

    createPieChart('pieChart', categories, totals, {
      cutout: '80%',
      colors: [
        "#FF3ECF","#FF6F00","#8605f2","#191cff",
        "#ffe500","#00FFFF","#ff194d","#52FF88"
      ],
      pluginsArray: [centerTextStatic]
    });

    // === CYCLING PIE ===
    let currentIndex = 0;
    const canvasId = "widgetPie";

    const centerTextCycle = {
      id: 'centerTextCycle',
      afterDraw(chart) {
        const { ctx, chartArea: { width, height } } = chart;
        ctx.save();
        ctx.font = "bold 18px Saira";
        ctx.fillStyle = "#FFFFFF";
        ctx.textAlign = "center";
        ctx.textBaseline = "middle";
        ctx.fillText(categories[currentIndex], width / 2, height / 2 - 20);
        ctx.font = "bold 28px Saira";
        ctx.fillText(`${((totals[currentIndex]/totalSum)*100).toFixed(1)}%`, width / 2, height / 2 + 10);
        ctx.restore();
      }
    };

    // initial chart
    const cycleChart = createPieChart(
      canvasId,
      [categories[0], "Other"],
      [totals[0], totalSum - totals[0]],
      {
        cutout: '90%',
        pluginsArray: [centerTextCycle],
        colors: ["#FF3ECF", "#2B2D30"]
      }
    );

    // update every 5s
    setInterval(() => {
      currentIndex = (currentIndex + 1) % categories.length;
      const value = totals[currentIndex];
      cycleChart.data.labels = [categories[currentIndex], "Other"];
      cycleChart.data.datasets[0].data = [value, totalSum - value];
      cycleChart.update();
    }, 5000);

  })
  .catch(err => console.error("Error fetching monthly totals:", err));
