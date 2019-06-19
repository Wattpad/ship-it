import React from 'react'
import ReactExpandableGrid from './components/TileGrid'
import './App.css'
import TopBar from './components/TopBar'
import axios from 'axios';
import { CircularProgress } from '@material-ui/core';

const urljoin = require('url-join')

const apiAddress = urljoin('https://' + window.location.hostname, 'releases') // Point to local IP for testing

class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {}
    axios.get(apiAddress).then(response => {
      console.log(response)
      this.setState({data : response.data })
    })
  }

  render() {
    return (
      <div className="App">
        <TopBar />
        {
          this.state.data ? 
          <ReactExpandableGrid
            gridData={this.state.data}
            detailHeight={300}
            ExpandedDetail_image_size={300}
            cellSize={250}
            apiAddress={apiAddress}
            ExpandedDetail_closeX_bool={false}
          /> : <CircularProgress />
        }
      </div>
    )
  }
}

export default App
