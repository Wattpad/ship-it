import React from 'react';
import ReactExpandableGrid from './components/TileGrid';
import './App.css';
import TopBar from './components/TopBar';

var data = [
  { // Formatted as it will be from the real API
    name: "homeslice",
    created: "2019-06-13T14:09:47.781282",
    lastDeployed: "2019-06-13T14:09:47.78Z",
    owner: {
      team: "squad-cd",
      slack: "#squad-cd"
    },
    autoDeploy: true,
    code: {
      github: "https://github.com/Wattpad/highlander/wattpad/src/services/homeslice",
      ref: "adf098ad00a8d76d5ad5ad4ada5dad4ad"
    },
    build: {
      travis: "" // Can get the build URL from GitHub Checks API
    },
    monitoring: {
      datadog: {
        dashboard: "https://app.datadoghq.com/dashboard/4px-qaj-tnc/home-v2?tile_size=m",
        monitors: "https://app.datadoghq.com/monitors/manage?q=home"
      },
      sumologic: "https://service.us2.sumologic.com/ui/#/search/KRBpz5OodF4HcTasdIEYkuhguVAEennEj7xIV8ke"
    },
    artifacts: {
      docker: {
        image: "723255503624.dkr.ecr.us-east-1.amazonaws.com/homeslice",
        tag: "adf098ad00a8d76d5ad5ad4ada5dad4ad"
      },
      chart: {
        path: "https://charts.wattpadhq.com/homeslice",
        version: "1.2.3"
      }
    },
    deployment: {
      status: "deploying", // deploying, deployed, rolled_back, failed
      // NOTE resources separate endpoint
    }
  }
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
        ExpandedDetail_closeX_bool={false}
      />
    </div>
  );
}

export default App;
