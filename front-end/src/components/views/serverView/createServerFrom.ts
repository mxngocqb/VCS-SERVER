import { ServerStatus } from "@/plugins/axios/server/interfaces";
import { serverService } from "@/plugins/axios/server/serverService";
import { useMainStore } from "@/stores/mainStore";
import { useForm } from "vee-validate";
import * as yup from "yup";

export default function () {
    const schema = yup.object({
        serverName: yup.string().required(),
        ipv4: yup
            .string()
            .matches(
                /^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$/,
                "This is not format of IPv4"
            )
            .required(),
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

    const [serverName, serverNameAttrs] = defineField("serverName");
    const [ipv4, ipv4Attrs] = defineField("ipv4");
    const [status, statusAttrs] = defineField("status");
    const mainStore = useMainStore();

    const onSubmit = handleSubmit(async (value) => {
        await serverService
            .createServer({
                name: value["serverName"],
                ip: value["ipv4"],
                status: value["status"] ? true : false,
            })
            .then((res) => {
                mainStore.showNofitication({
                    typeNotification: "infor",
                    title: "Create server",
                    content: "Create server successfuly",
                });
            });
        return true;
    });
    return {
        resetForm,
        setFieldValue,
        serverName,
        ipv4,
        status,
        onSubmit,
        errors,
        isValidating,
    };
}
