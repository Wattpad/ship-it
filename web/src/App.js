import React from 'react';
import ReactExpandableGrid from './components/TileGrid';
import './App.css';
import TopBar from './components/TopBar';

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
