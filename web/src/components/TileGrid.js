/* eslint react/prop-types: 0 */
/* eslint react/jsx-no-bind: 0 */

import React from 'react'
import PropTypes from 'prop-types'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'
import Chip from '@material-ui/core/Chip'
import DoneIcon from '@material-ui/icons/Done'
import ExpandIcon from '@material-ui/icons/ExpandMore'
import IconButton from '@material-ui/core/IconButton'
import SelectionDialog from './SelectionDialog'

import TimePassed from '../assets/time_passed.png'
import SlackIcon from '../assets/slack_icon.png'
import DockerIcon from '../assets/docker_icon.png'
import ExpandedDetail from './ExpandedDetail'

import { createMuiTheme, MuiThemeProvider } from '@material-ui/core/styles'

const deployTagTheme = createMuiTheme({
  palette: {
    primary: {
      main: '#4caf50'
    },
    secondary: {
      main: '#f44336'
    },
    default: {
      main: '#9e9e9e'
    }
  }
})

const imgAlt = "not found"

class SingleGridCell extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      expanded: false,
      selected_id: '',
      window_width: window.innerWidth,
      repoSelector: false,
    }
  }

  cellClick(event) {
    this.props.handleCellClick(event)
  }

  render() {
    var SingleGridCellStyle = {
      backgroundSize: this.props.cellSize,
      width: this.props.cellSize,
      height: this.props.cellSize,
      display: 'inline-block',
      margin: this.props.cellMargin,
      marginBottom: 25,
      position: 'relative'
    }

    var cardStyle = {
      width: this.props.cellSize,
      height: this.props.cellSize
    }

    // Re written to put material ui components in the tile original component only took images
    var deployDate = new Date(this.props.SingleGridCellData.lastDeployed)
    return (
      <div style={SingleGridCellStyle} id={this.props.id} className='SingleGridCell'>
        <div>
          <Card style={cardStyle}>
            <CardContent>
              <Typography variant="h5" component="h2">
                {this.props.SingleGridCellData.name}
              </Typography>
              <MuiThemeProvider theme={deployTagTheme}>
                <Chip
                  icon={<DoneIcon />}
                  label="Deployed"
                  color="primary"
                  variant="outlined"
                  clickable
                />
              </MuiThemeProvider>
              <div>
                <IconButton>
                  <img src={TimePassed} alt={imgAlt} />
                </IconButton>
                {deployDate.toDateString()}
              </div>
              <div className='row-align'>
                <SelectionDialog />
                <IconButton>
                  <img src={SlackIcon} width="32" height="32" alt={imgAlt} />
                </IconButton>
                <IconButton>
                  <img src={DockerIcon} width="32" height="32" alt={imgAlt} />
                </IconButton>
              </div>
              <div>
                <IconButton onClick={this.cellClick.bind(this)}>
                  <ExpandIcon />
                </IconButton>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    )
  }
}

var cardID
class ReactExpandableGrid extends React.Component {
  
  constructor(props) {
    super(props)

    this.state = {
      expanded: false,
      selected_id: '',
      gridData: this.props.gridData
    }
  }

  findIdNode(target) { // takes event target and returns the id (written to make sure the expanded detail renders from the correct reference point as in the original component each tile on the grid was just an img tag)
    var t = target
    var found = false
    while (!found && t != null) {
      if (String(t.className) === "SingleGridCell") {
        found = true
        return parseInt(t.id.substring(10))
      }
      t = t.parentNode
    }
    return -1
  }
  
  renderExpandedDetail(target) {
    var thisIdNumber = this.findIdNode(target)
    var detail = document.getElementById('expandedDetail')

    var ol = document.getElementById("grid_cell_" + thisIdNumber.toString()).parentNode
    var lengthOfList = parseInt(ol.childNodes.length)
    var startingIndex = thisIdNumber + 1

    var insertedFlag = false

    ol.insertBefore(detail, ol.childNodes[lengthOfList])

    for (var i = startingIndex; i < lengthOfList; i++) {
      if (ol.childNodes[i].className === 'SingleGridCell') {
        if (ol.childNodes[i].offsetTop !== ol.childNodes[thisIdNumber].offsetTop) {
          ol.childNodes[i].insertAdjacentElement('beforebegin', detail)
          insertedFlag = true
          break
        }
      }
    }

    if (insertedFlag === false) {
      ol.childNodes[lengthOfList - 1].insertAdjacentElement('afterend', detail)
    }
  }

  closeExpandedDetail() {
    this.setState({
      expanded: false,
      selected_id: ''
    }, function afterStateChange() {
      var detail = document.getElementById('expandedDetail')
      detail.style.display = 'none'
    })
  }

  handleCellClick(event) {
    var target = event.target
    cardID = this.findIdNode(target)
    if (this.state.expanded) {
      if (this.state.selected_id === event.target.id) { // Clicking on already opened detail
        this.closeExpandedDetail()
        this.renderExpandedDetail(target)
      } else { // Clicking on a different thumbnail, when detail is already expanded
        this.setState({
          expanded: true,
          selected_id: event.target.id
        }, function afterStateChange() {
          var detail = document.getElementById('expandedDetail')

          this.renderExpandedDetail(target)

          detail.style.display = 'block'
        })
      }
    } else {
      this.setState({
        expanded: true,
        selected_id: event.target.id
      }, function afterStateChange() {
        var detail = document.getElementById('expandedDetail')

        this.renderExpandedDetail(target)

        detail.style.display = 'block'
      })
    }
  }

  generateGrid() {
    var grid = []
    var idCounter = -1 // To help simplify mapping to object array indices. For example, <li> with 0th id corresponds to 0th child of <ol>
    var gridData = this.state.gridData
    for (var i in gridData) {
      idCounter = idCounter + 1
      var thisUniqueKey = 'grid_cell_' + idCounter.toString()
      grid.push(<SingleGridCell handleCellClick={this.handleCellClick.bind(this)} key={thisUniqueKey} id={thisUniqueKey} cellMargin={this.props.cellMargin} SingleGridCellData={gridData[i]} cellSize={this.props.cellSize} />)
    }

    var cssforExpandedDetail = {
      backgroundColor: this.props.detailBackgroundColor,
      height: this.props.detailHeight,
      display: 'none',
      position: 'relative',
      padding: '20px',
      transition: 'display 2s ease-in-out 0.5s',
    }

    grid.push( // Expanded Detail here
      <li style={cssforExpandedDetail} key='expandedDetail' id='expandedDetail'>
        <ExpandedDetail data={gridData} id={cardID}/>
      </li>
    )

    return grid
  }

  render() {
    var rows = this.generateGrid()

    var cssForGridDetailExpansion = {
      width: '100%',
      position: 'relative'
    }

    var cssForGridList = {
      listStyle: 'none',
      padding: 0,
      display: 'inline-block',
      width: '100%'
    }

    var cssForTheGridHolder = {
      width: '100%',
      backgroundColor: this.props.bgColor,
      margin: 0,
      textAlign: 'center'
    }
    return (
      <div id='GridDetailExpansion' style={cssForGridDetailExpansion}>
        <div id='theGridHolder' style={cssForTheGridHolder}>
          <ol id='gridList' style={cssForGridList}>
            {rows}
          </ol>
        </div>
      </div>
    )
  }
}

ReactExpandableGrid.propTypes = {
  gridData: PropTypes.array,
  cellSize: PropTypes.number,
  cellMargin: PropTypes.number,
  bgColor: PropTypes.string,
  detailHeight: PropTypes.number,
  detailBackgroundColor: PropTypes.string,
}

ReactExpandableGrid.defaultProps = {
  gridData: [],
  cellSize: 250,
  cellMargin: 25,
  bgColor: '#f2f2f2',
  detailHeight: 300,
  detailBackgroundColor: '#D9D9D9',
}

export default ReactExpandableGrid
