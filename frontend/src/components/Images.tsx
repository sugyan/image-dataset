import CSS from "csstype";
import H from "history";
import React, { useEffect, useState, useRef } from "react";
import { useHistory, useLocation } from "react-router";
import { Link as RouterLink, LinkProps } from "react-router-dom";
import {
    Button, Card, CardActionArea, CardMedia, CardContent,
    Box, CircularProgress, Grid, Link, Typography,
    Theme, makeStyles, createStyles,
} from "@material-ui/core";

import SearchBox from "./SearchBox";
import { ImageResponse } from "../common/interfaces";

const useStyles = makeStyles((theme: Theme) => {
    return createStyles({
        progress: {
            margin: theme.spacing(2),
        },
        card: {
            maxWidth: 144,
        },
        cardContent: {
            padding: theme.spacing(1.5),
        },
        media: {
            height: 144,
            width: 144,
        },
    });
});

const Images: React.FC = () => {
    const classes = useStyles();
    const history = useHistory();
    const location = useLocation();
    const last = useRef<string>();
    const [images, setImages] = useState<ImageResponse[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const loadImages = (history: H.History, location: H.Location, images: ImageResponse[]) => {
        const params = new URLSearchParams(location.search);
        params.set("count", "100");
        if (last.current) {
            params.set("id", last.current);
        }
        setLoading(true);
        fetch(`/api/images?${params}`).then((res: Response) => {
            if (res.ok) {
                return res.json();
            }
            if (res.status === 401) {
                history.push("/");
                return;
            }
            throw new Error(res.statusText);
        }).then((data: ImageResponse[]) => {
            if (data.length === 0) {
                return;
            }
            last.current = data[data.length - 1].id;
            const ids = new Set(images.map((value: ImageResponse) => value.id));
            setImages(images.concat(data.filter((value: ImageResponse) => {
                return !ids.has(value.id);
            })));
        }).catch((err: Error) => {
            window.console.error(err.message);
        }).finally(() => {
            setLoading(false);
        });
    };
    useEffect(() => {
        if (location.search.length === 0) {
            return;
        }
        if (images.length > 0) {
            if (last.current && last.current === images[images.length - 1].id) {
                return;
            }
        }
        loadImages(history, location, images);
    }, [history, location, images]);
    useEffect(() => {
        setImages([]);
    }, [location]);
    const cards = images.map((image: ImageResponse) => {
        const link = React.forwardRef<HTMLAnchorElement, Omit<LinkProps, "to">>(
            (props, ref) => {
                const to = {
                    pathname: `/image/${image.id}`,
                    search: location.search,
                };
                return (
                  <RouterLink innerRef={ref} to={to} {...props} />
                );
            },
        );
        const style: CSS.Properties = {
            borderStyle: "solid",
            borderWidth: "2px",
            borderColor: "transparent",
        };
        switch (image.status) {
        case 1:
            style.borderColor = "lightcoral";
            break;
        case 2:
            style.borderColor = "gray";
            break;
        case 3:
            style.borderColor = "lightgreen";
            break;
        }
        return (
          <Box key={image.id} m={0.2}>
            <Card className={classes.card} style={style}>
              <Link component={link}>
                <CardActionArea>
                  <CardMedia
                      className={classes.media}
                      image={image.image_url} />
                </CardActionArea>
              </Link>
              <CardContent className={classes.cardContent}>
                <Typography variant="body2">
                  {image.size}×{image.size}
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  {image.label_name}
                </Typography>
              </CardContent>
            </Card>
          </Box>
        );
    });
    const progress = (
      <Grid container justify="center">
        <CircularProgress className={classes.progress} />
      </Grid>
    );
    return (
      <React.Fragment>
        <h2>Images</h2>
        <SearchBox />
        <Box display="flex" flexWrap="wrap" mt={2}>
          {cards}
        </Box>
        {loading ? progress :<Grid container justify="center">
          <Box mt={2}>
            <Button color="inherit" onClick={() => loadImages(history, location, images)}>More</Button>
          </Box>
        </Grid>}
      </React.Fragment>
    );
};

export default Images;
