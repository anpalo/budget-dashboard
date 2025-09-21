async function displayTotalSavings() {
    try {
        const res = await fetch("/api/total-savings");
        const data = await res.json();
        const krw = Number(data.krw);

        const resRates = await fetch("/api/currency");
        const rates = await resRates.json();
        const usdRate = rates["USD"] || 0;
        const usd = krw * usdRate;

        const widget = document.getElementById("totalSavingsWidget");
        widget.innerHTML = `
            <div style="padding: 10px; display:flex; flex-direction:column; align-items:center; justify-content:center; height:100%; width:100%;">
                <h2 style="color:#0085ff; margin:0;">Total Savings</h2>
                <div style="font-size:1.5rem; font-weight:bold;">${krw.toLocaleString('ko-KR')}원</div>
                <div style="color:#555;">${usd.toLocaleString(undefined, {maximumFractionDigits:2})} USD</div>
            </div>
        `;
    } catch (err) {
        console.error("Error loading total savings:", err);
        document.getElementById("totalSavingsWidget").innerText = "Error loading savings.";
    }
}

displayTotalSavings();

function createPieChart(canvasId, labels, data, options = {}) {
    const ctx = document.getElementById(canvasId).getContext('2d');
    return new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                data: data,
                backgroundColor: labels.map((l,i) => options.colors?.[i] || "#2B2D30"),
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

fetch("/api/monthly-totals")
.then(res => res.json())
.then(data => {
    const month = "Jan";
    const categories = data[month].Totals.map(e => e.Category);
    const totals = data[month].Totals.map(e => e.Total);

    const centerTextStatic = {
        id: "centerTextStatic",
        afterDraw(chart) {
            const { ctx, chartArea: { width, height } } = chart;
            ctx.save();
            ctx.font = "bold 18px Saira";
            ctx.fillStyle = "#FFFFFF";
            ctx.textAlign = "center";
            ctx.textBaseline = "middle";
            ctx.fillText("Spending in " + month, width / 2, 20);
            ctx.restore();
        }
    };

    createPieChart('pieChart', categories, totals, {
        cutout: '80%',
        colors: [
            "#FF3ECF",
            "#FF6F00",
            "#8605f2",
            "#191cff",
            "#ffe500",
            "#00FFFF",
            "#ff194d",
            "#52FF88"
        ],
        pluginsArray: [centerTextStatic]
    });

    const highlightCategory = "Development";
    const highlightIndex = categories.indexOf(highlightCategory);
    const highlightValue = totals[highlightIndex];
    const otherValue = totals.reduce((sum, val, i) => i !== highlightIndex ? sum + val : sum, 0);

    createPieChart('widgetPie', [highlightCategory, "Other"], [highlightValue, otherValue], {
        cutout: '90%',
        colors: ["#FF3ECF", "#2B2D30"]
    });
});
