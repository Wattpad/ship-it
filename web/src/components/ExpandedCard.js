import React from 'react'
import { Dialog, DialogContent, DialogActions, List, Button, DialogTitle, Typography, ListItem, ListItemText } from '@material-ui/core';

class ExpandedCard extends React.Component {
    constructor(props) {
        super(props)
        this.state = {}
    }

    render() {
        return (
            <Dialog open={this.props.open} onClose={this.props.handleClose}>
                <DialogTitle>Hello World</DialogTitle>
            </Dialog>
        )
    }
}

export default ExpandedCard