import { ServerStatus } from "@/plugins/axios/server/interfaces";
import { useServerStore } from "@/stores/serverStore";
import { useForm } from "vee-validate";
import * as yup from "yup";

export default function () {
    const schema = yup.object({
        createdAtFrom: yup.date().optional(),
        createdAtTo: yup.date().min(yup.ref("createdAtFrom")).optional(),
        updatedAtFrom: yup.date().optional(),
        updatedAtTo: yup.date().min(yup.ref("createdAtFrom")).optional(),
        status: yup.boolean(),
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

    const [createdAtFrom, createdAtFromAttrs] = defineField("createdAtFrom");
    const [createdAtTo, createdAtToAttrs] = defineField("createdAtTo");
    const [updatedAtFrom, updatedAtFromAttrs] = defineField("updatedAtFrom");
    const [updatedAtTo, updatedAtToAttrs] = defineField("updatedAtTo");
    const [status, statusAttrs] = defineField("status");
    const serverStore = useServerStore();

    const onSubmit = handleSubmit(async (value) => {
        await serverStore.updateFilterServer({
            createdAtFrom: value["createdAtFrom"]
                ? new Date(value["createdAtFrom"]).toISOString()
                : undefined,
            createdAtTo: value["createdAtTo"]
                ? new Date(value["createdAtTo"]).toISOString()
                : undefined,
            updatedAtFrom: value["updatedAtFrom"]
                ? new Date(value["updatedAtFrom"]).toISOString()
                : undefined,
            updatedAtTo: value["updatedAtFrom"]
                ? new Date(value["updatedAtTo"]).toISOString()
                : undefined,
            status: status.value
                ? ServerStatus.ON
                : status.value === undefined
                ? undefined
                : ServerStatus.OFF,
        });
        return true;
    });
    return {
        resetForm,
        setFieldValue,
        onSubmit,
        errors,
        isValidating,
        status,
        createdAtFrom,
        createdAtTo,
        updatedAtFrom,
        updatedAtTo,
    };
}
