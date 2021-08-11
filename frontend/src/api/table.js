import { instance as axios } from "../plugins/axios";

export const table = {
  async dropTable(table) {
    try {
      return await axios.delete("/table/drop/" + table);
    } catch (err) {
      return err;
    }
  },
};
