import { ServerStatus, TypeSort } from "@/plugins/axios/server/interfaces";
import { serverService } from "@/plugins/axios/server/serverService";
import { useServerStore } from "@/stores/serverStore";
import { useUserStore } from "@/stores/userStore";
import { useForm } from "vee-validate";
import * as yup from "yup";

export default function () {
    const schema = yup.object({
        status: yup.boolean(),
        fromPage: yup
            .number()
            .min(1)
            .transform((_, val) => (val === "" ? undefined : Number(val)))
            .optional(),
        toPage: yup
            .number()
            .min(yup.ref("fromPage"))
            .transform((_, val) => (val === "" ? undefined : Number(val)))
            .optional(),
        sort: yup.string().optional(),
        sortBy: yup.string().optional(),
        oder: yup.string().optional(),
    });

    const {
        handleSubmit,
        errors,
        isValidating,
        defineField,
        setFieldValue,
        resetForm,
    } = useForm({
        validationSchema: schema,
    });

  
    const [status, statusAttrs] = defineField("status");
    const [fileName, fileNameAttrs] = defineField("fileName");
    const [fromPage, pageAttrs] = defineField("fromPage");
    const [toPage, toPageAttrs] = defineField("toPage");
    const [sort, sortAttrs] = defineField("sort");
    const [sortBy, sortByAttrs] = defineField("sortBy");
    const [pageSize, pageSizeAttrs] = defineField("pageSize");
    const [order, orderAttrs] = defineField("order");

    const serverStore = useServerStore();

    const userStore = useUserStore();

    const onSubmit = handleSubmit(async (value) => {
        console.log(value);
        
        serverService.exportServer({
            limit: value.pageSize,
            offset: value.fromPage,
            status: value.status ? "true" : "false",
            field: value.sortBy,
            order: value.sort,
        });
        
        return true;
    });
    return {
        resetForm,
        setFieldValue,
        onSubmit,
        errors,
        pageSize,
        isValidating,
        status,
        fileName,
        fromPage,
        toPage,
        sort,
        sortBy,
    };
}
