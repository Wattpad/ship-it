import React from 'react'
import ReactExpandableGrid from './components/TileGrid'
import './App.css'
import TopBar from './components/TopBar'
import axios from 'axios';
import { CircularProgress } from '@material-ui/core';
import urljoin from 'url-join'

//const API_ADDRESS = urljoin('https://' + window.location.hostname, 'api') // Point to local IP for testing
const API_ADDRESS = 'http://localhost:8080/'
class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {}
    axios.get(urljoin(API_ADDRESS, 'releases')).then(response => {
      this.setState({data : response.data})
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
            API_ADDRESS={API_ADDRESS}
            ExpandedDetail_closeX_bool={false}
          /> : <CircularProgress />
        }
      </div>
    )
  }
}

export default App
