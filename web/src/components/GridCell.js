import React from 'react'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'
import ExpandIcon from '@material-ui/icons/ExpandMore'
import IconButton from '@material-ui/core/IconButton'
import SelectionDialog from './SelectionDialog'

import TimePassed from '../assets/time_passed.png'
import SlackIcon from '../assets/slack_icon.png'
import DockerIcon from '../assets/docker_icon.png'
import StatusChip from './StatusChip';

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
                            <StatusChip status={this.props.SingleGridCellData.deployment.status} />
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

export default SingleGridCell