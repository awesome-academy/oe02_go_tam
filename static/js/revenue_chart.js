document.addEventListener("DOMContentLoaded", function () {
    const labels = window.revenueLabels || [];
    const data = window.revenueData || [];

    const ctx = document.getElementById("revenueChart").getContext("2d");
    new Chart(ctx, {
        type: "bar",
        data: {
            labels: labels.map((m) => "" + m),
            datasets: [{
                label: "Revenue",
                data: data,
                backgroundColor: "rgba(59, 130, 246, 0.6)",
                borderColor: "rgba(59, 130, 246, 1)",
                borderWidth: 1,
            }],
        },
        options: {
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function (value) {
                            return "$" + value;
                        },
                    },
                },
            },
        },
    });
});
