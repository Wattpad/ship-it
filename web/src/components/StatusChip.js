import React from 'react'
import Chip from '@material-ui/core/Chip'
import DoneIcon from '@material-ui/icons/Done'
import RollBackIcon from '@material-ui/icons/Cached'
import FailIcon from '@material-ui/icons/Clear'

class StatusChip extends React.Component {

    getChip(state) {
        switch(state) {
            case 'deploying':
                return <Chip
                    icon={<DoneIcon />}
                    label="Deployed"
                    color="primary"
                    variant="outlined"
                    clickable
                />
            case 'failed':
                return <Chip
                    icon={<FailIcon />}
                    label="Failed"
                    color="primary"
                    variant="outlined"
                    clickable
                />
            case 'rollback':
                return <Chip
                    icon={<RollBackIcon />}
                    label="Rolled Back"
                    color="primary"
                    variant="outlined"
                    clickable
                />
            default:
                return null
        }
    }

    render() {
        return (
            <div>
                {console.log(this.props.status)}
                {this.getChip(this.props.status)}
            </div>
        )
    }
}

export default StatusChip