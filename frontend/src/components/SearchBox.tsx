import React, { useState, ChangeEvent, useEffect, useRef } from "react";
import { useHistory, useLocation } from "react-router";
import {
    Box, Button, Collapse, MenuItem, Paper,
    List, ListItem, ListItemIcon, ListItemText,
    FormControl, FormControlLabel, FormLabel, InputLabel,
    Select, TextField, Radio, RadioGroup,
    Theme, makeStyles, createStyles,
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
    const nameValue = useRef<string>();
    const params = new URLSearchParams(location.search);
    const [name, setName] = useState<string>(params.get("name") || "");
    const [status, setStatus] = useState<string>(params.get("status") || "all");
    const [size, setSize] = useState<string>(params.get("size") || "all");
    const [sort, setSort] = useState<string>(params.get("sort") || "id");
    const [order, setOrder] = useState<string>(params.get("order") || "asc");
    const [expand, setExpand] = useState<boolean>(false);
    const onChangeName = (e: ChangeEvent<{ name?: string, value: any }>) => {
        setName(e.target.value);
        const value = e.target.value;
        nameValue.current = value;
        setTimeout(() => {
            if (nameValue.current === value) {
                const params = new URLSearchParams(location.search);
                params.set("name", value);
                history.replace({
                    pathname: history.location.pathname,
                    search: params.toString(),
                });
            }
        }, 500);
    };
    const onChangeStatus = (e: ChangeEvent<{ name?: string, value: any }>) => {
        setStatus(e.target.value);
    };
    const onChangeSize = (e: ChangeEvent<{ name?: string; value: any }>) => {
        setSize(e.target.value);
    };
    const onChangeSort = (event: ChangeEvent<HTMLInputElement>) => {
        setSort(event.target.value);
    };
    const onChangeOrder = (event: ChangeEvent<HTMLInputElement>) => {
        setOrder(event.target.value);
    };
    const resetForm = () => {
        setStatus("all");
        setSize("all");
        setSort("id");
        setOrder("asc");
    };
    useEffect(() => {
        const params = new URLSearchParams({ status, size, sort, order });
        const name = new URLSearchParams(location.search).get("name");
        if (name) {
            params.set("name", name);
        }
        if (`?${params}` !== location.search) {
            history.replace({
                pathname: history.location.pathname,
                search: params.toString(),
            });
        }
    }, [status, size, sort, order, history, location]);
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
              <TextField
                  label="Name"
                  InputLabelProps={{ shrink: true }}
                  value={name}
                  onChange={onChangeName} />
            </FormControl>
            <FormControl className={classes.formControl}>
              <InputLabel>
                Status
              </InputLabel>
              <Select value={status} onChange={onChangeStatus}>
                <MenuItem value={"all"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">
                    -----
                  </Box>
                </MenuItem>
                <MenuItem value={"0"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">
                    = Ready
                  </Box>
                </MenuItem>
                <MenuItem value={"1"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">
                    = NG
                  </Box>
                </MenuItem>
                <MenuItem value={"2"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">
                    = Pending
                  </Box>
                </MenuItem>
                <MenuItem value={"3"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">
                    = OK
                  </Box>
                </MenuItem>
                <MenuItem value={"4"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">
                    = Predicted
                  </Box>
                </MenuItem>
              </Select>
            </FormControl>
            <FormControl className={classes.formControl}>
              <InputLabel>
                Size
              </InputLabel>
              <Select value={size} onChange={onChangeSize}>
                <MenuItem value={"all"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">
                    -----
                  </Box>
                </MenuItem>
                <MenuItem value={"256"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">
                    &gt;= 256
                  </Box>
                </MenuItem>
                <MenuItem value={"512"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">
                    &gt;= 512
                  </Box>
                </MenuItem>
                <MenuItem value={"1024"}>
                  <Box fontSize="body1.fontSize" fontFamily="Monospace">
                    &gt;= 1024
                  </Box>
                </MenuItem>
              </Select>
            </FormControl>
            <FormControl component="fieldset" className={classes.formControl}>
              <FormLabel component="legend">
                Order by
              </FormLabel>
              <RadioGroup
                  aria-label="sort"
                  name="sort"
                  value={sort}
                  onChange={onChangeSort}
              >
                <FormControlLabel
                    value="id"
                    control={<Radio />}
                    label="ID"
                />
                <FormControlLabel
                    value="updated_at"
                    control={<Radio />}
                    label="Updated At"
                />
                <FormControlLabel
                    value="published_at"
                    control={<Radio />}
                    label="Published At"
                />
              </RadioGroup>
            </FormControl>
            <FormControl component="fieldset" className={classes.formControl}>
              <FormLabel component="legend">
                &nbsp;
              </FormLabel>
              <RadioGroup
                  aria-label="order"
                  name="order"
                  value={order}
                  onChange={onChangeOrder}
              >
                <FormControlLabel
                    value="asc"
                    control={<Radio />}
                    label="Asc"
                />
                <FormControlLabel
                    value="desc"
                    control={<Radio />}
                    label="Desc"
                />
              </RadioGroup>
            </FormControl>
            <FormControl component="fieldset" className={classes.formControl}>
              <FormLabel component="legend">
                &nbsp;
              </FormLabel>
              <Button color="primary" onClick={() => resetForm()}>Reset</Button>
            </FormControl>
          </Box>
        </Collapse>
      </Paper>
    );
};

export default SearchBox;
