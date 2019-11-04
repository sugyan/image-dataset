import React, { useState, ChangeEvent, useEffect } from "react";
import { useHistory, useLocation } from "react-router";
import {
    Paper, Box, Collapse,
    List, ListItem, ListItemIcon, ListItemText,
    FormControl, FormLabel, FormControlLabel,
    InputLabel, Select, MenuItem, RadioGroup, Radio,
    makeStyles, createStyles, Theme,
} from "@material-ui/core";
import { ExpandMore, ImageSearch, ExpandLess } from "@material-ui/icons";

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        formControl: {
            margin: theme.spacing(1),
            minWidth: 120,
        },
    }),
);

const SearchBox: React.FC = () => {
    const classes = useStyles();
    const history = useHistory();
    const location = useLocation();
    const params = new URLSearchParams(location.search);
    const [size, setSize] = useState<string>(params.get("size") || "all");
    const [sort, setSort] = useState<string>(params.get("sort") || "id");
    const [order, setOrder] = useState<string>(params.get("order") || "asc");
    const [expand, setExpand] = useState<boolean>(false);
    const onChangeSize = (e: ChangeEvent<{ name?: string; value: any }>) => {
        setSize(e.target.value);
    };
    const onChangeSort = (event: ChangeEvent<HTMLInputElement>) => {
        setSort(event.target.value);
    };
    const onChangeOrder = (event: ChangeEvent<HTMLInputElement>) => {
        setOrder(event.target.value);
    };
    useEffect(() => {
        const params = new URLSearchParams({ size, sort, order });
        if (`?${params}` !== location.search) {
            history.push({
                pathname: history.location.pathname,
                search: params.toString(),
            });
        }
    }, [size, sort, order, history, location]);
    return (
      <Paper square={true}>
        <List>
          <ListItem button onClick={() => setExpand(!expand)}>
            <ListItemIcon>
              <ImageSearch />
            </ListItemIcon>
            <ListItemText primary="Search" />
            {expand ? <ExpandLess /> : <ExpandMore />}
          </ListItem>
        </List>
        <Collapse in={expand} mountOnEnter unmountOnExit>
          <Box m={3}>
            <FormControl className={classes.formControl}>
              <InputLabel>Size</InputLabel>
              <Select value={size} onChange={onChangeSize}>
                <MenuItem value={"all"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">-----</Box>
                </MenuItem>
                <MenuItem value={"256"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">&gt;= 256</Box>
                </MenuItem>
                <MenuItem value={"512"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">&gt;= 512</Box>
                </MenuItem>
                <MenuItem value={"1024"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">&gt;= 1024</Box>
                </MenuItem>
              </Select>
            </FormControl>
            <FormControl component="fieldset" className={classes.formControl}>
              <FormLabel component="legend">Sort by</FormLabel>
              <RadioGroup aria-label="sort" name="sort" value={sort} onChange={onChangeSort}>
                <FormControlLabel value="id" control={<Radio />} label="ID" />
                <FormControlLabel value="posted_at" control={<Radio />} label="Posted At" />
              </RadioGroup>
            </FormControl>
            <FormControl component="fieldset" className={classes.formControl}>
              <FormLabel component="legend">Order</FormLabel>
              <RadioGroup aria-label="order" name="order" value={order} onChange={onChangeOrder}>
                <FormControlLabel value="asc" control={<Radio />} label="Asc" />
                <FormControlLabel value="desc" control={<Radio />} label="Desc" />
              </RadioGroup>
            </FormControl>
          </Box>
        </Collapse>
      </Paper>
    );
};

export default SearchBox;
