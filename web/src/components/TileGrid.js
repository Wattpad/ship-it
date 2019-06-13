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
import CompressIcon from '@material-ui/icons/ExpandLess'
import IconButton from '@material-ui/core/IconButton'
import Green from '@material-ui/core/colors/green';
import SelectionDialog from './SelectionDialog'
import CircularProgress from '@material-ui/core/CircularProgress'

import TimePassed from '../assets/time_passed.png'
import SlackIcon from '../assets/slack_icon.png'
import GitIcon from '../assets/octocat.png'
import DockerIcon from '../assets/docker_icon.png'
import KubeIcon from '../assets/kubernetes_icon.png'
import JenkinsIcon from '../assets/jenkins_icon.png'
import TravisIcon from '../assets/travis_icon.png'
import ExpandedDetail from './ExpandedDetail';

import { createMuiTheme, MuiThemeProvider } from '@material-ui/core/styles';

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

var keys = []

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
    console.log(keys)
    var cardStyle = {
      width: this.props.cellSize,
      height: this.props.cellSize
    }

    // Re written to put material ui components in the tile original component only took images
    return (
      <div style={SingleGridCellStyle} id={this.props.id} className='SingleGridCell'>
        <div>
          <Card style={cardStyle}>
            <CardContent >
              <Typography variant="h5" component="h2">
                ServiceName
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
                  <img src={TimePassed} />
                </IconButton>
                5 min ago
              </div>
              <div className='row-align'>
                <SelectionDialog />
                <IconButton>
                  <img src={SlackIcon} width="32" height="32" />
                </IconButton>
                <IconButton>
                  <img src={DockerIcon} width="32" height="32" />
                </IconButton>
              </div>
              <div>
                {
                  this.state.expanded // Not working
                    ?
                    <IconButton onClick={this.cellClick.bind(this)}>
                      <CompressIcon />
                    </IconButton>
                    :
                    <IconButton onClick={this.cellClick.bind(this)}>
                      <ExpandIcon />
                    </IconButton>
                }
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    )
  }
}

class ReactExpandableGrid extends React.Component {

  constructor(props) {
    super(props)

    this.state = {
      expanded: false,
      selected_id: '',
      gridData: JSON.parse(this.props.gridData)
    }
  }

  findIdNode(target) { // takes event target and returns the id (written to make sure the expanded detail renders from the correct reference point as in the original component each tile on the grid was just an img tag)
    var t = target
    var found = false
    while (!found && t != null) {
      console.log(t.className)
      if (String(t.className) === "SingleGridCell") {
        found = true;
        return parseInt(t.id.substring(10))
      }
      t = t.parentNode
    }
    return -1;
  }

  renderExpandedDetail(target) {
    var thisIdNumber = this.findIdNode(target)
    var detail = document.getElementById('expandedDetail')

    console.log(detail)
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
    var thisIdNumber = this.findIdNode(target)
    var className = event.target.className
    console.log(event)
    console.log(thisIdNumber);
    if (this.state.expanded) { // expanded == true
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
    } else { // expanded == false
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
    var k = []
    for (var i in gridData) {
      idCounter = idCounter + 1
      var thisUniqueKey = 'grid_cell_' + idCounter.toString()
      k.push(thisUniqueKey)
      grid.push(<SingleGridCell handleCellClick={this.handleCellClick.bind(this)} key={thisUniqueKey} id={thisUniqueKey} cellMargin={this.props.cellMargin} SingleGridCellData={gridData[i]} cellSize={this.props.cellSize} />)
    }
    keys = k
    var cssforExpandedDetail = {
      backgroundColor: this.props.detailBackgroundColor,
      height: this.props.detailHeight,
      width: this.props.detailWidth,
      display: 'none',
      position: 'relative',
      padding: '20px',
      transition: 'display 2s ease-in-out 0.5s',
    }

    var cssforExpandedDetailLeft
    var cssforExpandedDetailRight

    cssforExpandedDetailLeft = {
      width: this.props.ExpandedDetail_left_width,
      height: '100%',
      float: 'left',
      position: 'relative'
    }

    cssforExpandedDetailRight = {
      width: this.props.ExpandedDetail_right_width,
      height: '100%',
      float: 'right',
      position: 'relative'
    }

    // Make Mobile Friendly
    if (window.innerWidth < this.props.show_mobile_style_from_width) {
      cssforExpandedDetailLeft = {
        width: '0%',
        height: '100%',
        float: 'left',
        position: 'relative',
        display: 'none'
      }

      cssforExpandedDetailRight = {
        width: '100%',
        height: '100%',
        float: 'right',
        position: 'relative'
      }
    }

    var closeX
    if (this.props.ExpandedDetail_closeX_bool) {
      closeX = 'X'
    } else {
      closeX = ''
    }

    grid.push( // Expanded Detail here
      <li style={cssforExpandedDetail} key='expandedDetail' id='expandedDetail'>
        <ExpandedDetail />
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
      display: 'inline-block'
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
  gridData: PropTypes.string,
  cellSize: PropTypes.number,
  cellMargin: PropTypes.number,
  bgColor: PropTypes.string,
  detailWidth: PropTypes.string, // in %
  detailHeight: PropTypes.number,
  detailBackgroundColor: PropTypes.string,
  ExpandedDetail_closeX_bool: PropTypes.bool,
  show_mobile_style_from_width: PropTypes.number
}

ReactExpandableGrid.defaultProps = {
  gridData: "",
  cellSize: 250,
  cellMargin: 25,
  bgColor: '#f2f2f2',
  detailWidth: '100%',
  detailHeight: 300,
  detailBackgroundColor: '#D9D9D9',
  ExpandedDetail_closeX_bool: true,
  show_mobile_style_from_width: 600,
}

export default ReactExpandableGrid
