import React from 'react'
import ReactExpandableGrid from './components/TileGrid'
import './App.css'
import TopBar from './components/TopBar'
import axios from 'axios';
import { CircularProgress } from '@material-ui/core';
import urljoin from 'url-join'

const API_ADDRESS = urljoin(window.location.protocol + '//' + window.location.host, 'api')
class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = { query: "" }
    axios.get(urljoin(API_ADDRESS, 'releases')).then(response => {
      this.setState({data : response.data})
    })
  }

  render() {
    return (
      <div className="App">
        <TopBar onQueryChanged={ q => this.setState({ query: q }) }/>
        {
          this.state.data ?
          <ReactExpandableGrid
            gridData={this.state.data}
            detailHeight={300}
            ExpandedDetail_image_size={300}
            cellSize={250}
            API_ADDRESS={API_ADDRESS}
            ExpandedDetail_closeX_bool={false}
            query={this.state.query}
          /> : <CircularProgress />
        }
      </div>
    )
  }
}

export default App
