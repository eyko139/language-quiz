import { addWord } from "@/app/lib/actions";
import { Button, TextField } from "@mui/material";


export default function Addword() {
    return (
        <form action={addWord}>
            <TextField name='word' id="outlined-basic" label="Outlined"
                variant="outlined" />
            <Button type='submit' color='success' variant="contained">Add Word</Button>
        </form>
    )
}
