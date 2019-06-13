import React from 'react';
import ReactExpandableGrid from './components/TileGrid';
import './App.css';
import TopBar from './components/TopBar';

var data = [
  {
    serviceName: "loki",
    deploymentStatus: "deployed",
    timeStamp: "5 min ago",
    chartVersion: "0.1.1",
    autoDeploy: true,
    mirandaPR: "https://github.com",
    highlanderPR: "https://github.com",
    slackLink: "https://slack.com",
    dockerLink: "https://www.docker.com/",
    dataDogDashboard: "dd",
    dataDogMonitor: "",
    sumoLogs: "",
    travisBuild: "",
    kubeResources: [
      "pod 1",
      "pod 2",
      "pod 3",
      "pod 4",
    ]
  }, {
    serviceName: "hermes",
    deploymentStatus: "deployed",
    timeStamp: "5 min ago",
    chartVersion: "0.1.1",
    autoDeploy: true,
    mirandaPR: "https://github.com",
    highlanderPR: "https://github.com",
    slackLink: "https://slack.com",
    dockerLink: "https://www.docker.com/",
    dataDogDashboard: "dd",
    dataDogMonitor: "",
    sumoLogs: "",
    travisBuild: "",
    kubeResources: [
      "pod 1",
      "pod 2",
      "pod 3",
      "pod 4",
    ]
  },
]

function App() {
  var dataString = JSON.stringify(data)
  return (
    <div className="App">
      <TopBar />
      <ReactExpandableGrid
        gridData={dataString}
        detailHeight={300}
        ExpandedDetail_image_size={300}
        cellSize={250}
        detailWidth='100%'
        ExpandedDetail_closeX_bool={false}
      />
    </div>
  );
}

export default App;
