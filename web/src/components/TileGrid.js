/* eslint react/prop-types: 0 */
/* eslint react/jsx-no-bind: 0 */

import React from 'react'
import PropTypes from 'prop-types'
import ExpandedDetail from './ExpandedDetail'
import SingleGridCell from './GridCell'

var cardID = 0
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
        <ExpandedDetail data={gridData[cardID]} />
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
