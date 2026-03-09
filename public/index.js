
moment.tz('Europe/Berlin')
Chart.defaults.font.size = 24;

let humChart, pwrChart

function toMoment(time) {


    let m = moment(time);
    m = m.subtract({hours: 1});
    return m
}

async function fetchTags() {
    let res = await fetch("/tags");
    let tags = await res.json();
    return tags;
}

async function fetchChartData(series, tag, from) {

    const queryParams = new URLSearchParams({tag, from});
    const url = `/${series}?${queryParams.toString()}`
    let res = await fetch(url);
    let data = await res.json();

    return data;
}

async function fetchTagData(tag, from) {

    let humidityDataPromise = fetchChartData("humidity",tag, from);
    let powerDataPromise =  fetchChartData("power",tag, from);
    let stateDataPromise = fetchChartData("state",tag, from);

    let [humidityData, powerData, stateData] = await Promise.all([humidityDataPromise, powerDataPromise,stateDataPromise])
    return {
        humidityData, powerData, stateData
    }

}

function initHumChart({humidityData, powerData, stateData}) {

    // let humidityData = await fetchChartData("humidity",tag, new Date().getTime());
    // let powerData = await fetchChartData("power",tag, new Date().getTime());
    // let stateData = await fetchChartData("state",tag, new Date().getTime());

    let humidityDataSet =  {
        label: 'Luftfeuchtigkeit',
        yAxisID: 'humidity',
        data: humidityData.map(v => ({x: toMoment(v.Time), y: v.Humidity})),
        backgroundColor: 'rgba(8,62,236,0.2)',
        borderColor: 'rgb(101,119,234)',
        borderWidth: 1
    }


    let powerDataSet =  {
        label: 'Apparent Power',
        yAxisID: 'power',
        data: powerData.map(v => ({x: toMoment(v.Time), y: v.ApparentPower})),
        backgroundColor: 'rgba(193,151,91,0.2)',
        borderColor: 'rgb(207,179,106)',
        borderWidth: 1
    }

    // let stateDataSet =  {
    //     label: 'An/Aus',
    //     data: stateData.map(v => ({x: v.Time, y: v.SwitchState})),
    //     backgroundColor: 'rgba(125,158,124,0.2)',
    //     borderColor: 'rgb(8,255,0)',
    //     borderWidth: 1
    // }

    chart = new Chart(document.getElementById('humChart').getContext('2d'), {
        type: 'line',
        data: {
            datasets: [humidityDataSet, powerDataSet]
        },
        options: {
            scales:{
                x:{
                    type: 'time',
                    ticks: {
                        autoSkip: true,
                        maxTicksLimit: 20
                    },
                    time: {
                        displayFormats: {minute: 'HH:mm'}
                    }
                },

                humidity: {
                    type: 'linear',
                    position: 'left',
                    ticks:
                        {
                            beginAtZero: true,
                        },
                    grid: { display: false }
                },
                power: {
                    type: 'linear',
                    position: 'right',
                    ticks:
                        {
                            beginAtZero: true,
                        },
                    grid: { display: false }
                },
            }
        }
    });
    return chart
}

function initPwrChart({humidityData, powerData, stateData}) {

    //
    // let humidityDataSet =  {
    //     label: 'Luftfeuchtigkeit',
    //     yAxisID: 'humidity',
    //     data: humidityData.map(v => ({x: toMoment(v.Time), y: v.Humidity})),
    //     backgroundColor: 'rgba(8,62,236,0.2)',
    //     borderColor: 'rgb(101,119,234)',
    //     borderWidth: 1
    // }


    let totalDataset =  {
        label: 'Total KWh',
        yAxisID: 'total',
        data: powerData.map(v => ({x: toMoment(v.Time), y: v.Total})),
        backgroundColor: 'rgba(193,151,91,0.2)',
        borderColor: 'rgb(207,179,106)',
        borderWidth: 1
    }

    let yesterdayDataset =  {
        label: 'Gestern KWh',
        yAxisID: 'yesterday',
        data: powerData.map(v => ({x: toMoment(v.Time), y: v.Yesterday})),
        backgroundColor: 'rgba(110,70,70,0.2)',
        borderColor: 'rgb(207,106,114)',
        borderWidth: 1
    }



    chart = new Chart(document.getElementById('pwrChart').getContext('2d'), {
        type: 'line',
        data: {
            datasets: [totalDataset, yesterdayDataset]
        },
        options: {
            scales:{
                x:{
                    type: 'time',
                    ticks: {
                        autoSkip: true,
                        maxTicksLimit: 20
                    },
                    time: {
                        displayFormats: {minute: 'HH:mm'}
                    }
                },

                total: {
                    type: 'linear',
                    position: 'left',
                    ticks:
                        {
                            beginAtZero: true,
                        },
                    grid: { display: false }
                },
                yesterday: {
                    type: 'linear',
                    position: 'right',
                    ticks:
                        {
                            beginAtZero: true,
                        },
                    grid: { display: false }
                },
            }
        }
    });
    return chart
}

async function  loadChartData(tag, date) {
    console.log(date)
    let data = await fetchTagData(tag, date);

    humChart = initHumChart(data);
    pwrChart = initPwrChart(data);
}


async function populateOptions() {

    let tagSelectElement = document.getElementById("tag-select")
    tagSelectElement.active = false
    let tags = await fetchTags();
    tagSelectElement.innerHTML = '';
    for (let tag of tags) {
        let option = document.createElement("option");
        option.value = tag;
        option.text = tag;
        tagSelectElement.appendChild(option);
    }
    tagSelectElement.active = true
    let rangeSelectElement = document.getElementById("range-select")
    document.getElementById("submit-select").addEventListener("click", function(event){
        event.preventDefault();
        console.log("submit clicked", rangeSelectElement.value, tagSelectElement.value);
        let m = moment().subtract({hours: rangeSelectElement.value});
        console.log(m);
        var fromDate = m.toDate();
        console.log(fromDate.toISOString());
        if (humChart !== undefined) {
            humChart.destroy();
        }
        if (pwrChart !== undefined) {
            pwrChart.destroy();
        }

        loadChartData(tagSelectElement.value, fromDate.toISOString());
    });


}



populateOptions();

let m = moment().subtract({hours: 3});
var fromDate = m.toDate();

loadChartData("Bad Unten", fromDate.toISOString());